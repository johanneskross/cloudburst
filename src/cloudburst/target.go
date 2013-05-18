package cloudburst

import (
	"container/list"
	"fmt"
	"strconv"
	"time"
	"times"
)

const toNanoseconds = 1000000000

type Target struct {
	Name          string
	Agents        list.List
	AgentChannel  chan bool
	Configuration TargetConfiguration
}

func NewTarget(targetConfiguration TargetConfiguration) *Target {
	agents := *list.New()
	channelSize := calcChannelSize(targetConfiguration.TimeSeries.Elements)
	agentChannel := make(chan bool, channelSize)
	return &Target{targetConfiguration.Name, agents, agentChannel, targetConfiguration}
}

func calcChannelSize(elements []*times.Element) int {
	channelSize := 0
	for i := 0; i < len(elements); i++ {
		value := int(elements[i].Value)
		if value > channelSize {
			channelSize = value
		}
	}
	return channelSize
}

func (t *Target) RunTimeSeries(c chan bool) {
	fmt.Printf("Running time series on target: %v\n", t.Name)
	
	startTime := time.Now().UnixNano() + int64(t.Configuration.RampUp) * toNanoseconds
	duration := t.Configuration.Duration
	
	for i := 0; i < duration; i++ {
		// wait until next interval is due
		nextInterval := (t.Configuration.TimeSeries.Elements[i].Timestamp * toNanoseconds) + startTime
		t.Wait(nextInterval)

		runningAgents := len(t.AgentChannel)
		runningNextAgents := int(t.Configuration.TimeSeries.Elements[i].Value)
		fmt.Printf("Update amount of agents to %v on target: %v\n", runningNextAgents, t.Name)
		
		// update amount of agents for this interval
		switch {
		case runningAgents < runningNextAgents:
			addAgents := runningNextAgents - runningAgents
			startAgents(t, addAgents)
		case runningAgents > runningNextAgents:
			reduceAgents := runningAgents - runningNextAgents
			interruptAgents(t, reduceAgents)
		}
	}
	c <- true
}

func (t *Target) Wait(nextInterval int64) {
	currentTime := time.Now().UnixNano()
	deltaTime := nextInterval - currentTime
	fmt.Printf("Target %v waits %v seconds for next interval\n", t.Name, deltaTime / toNanoseconds)
	if deltaTime > 0 {
		time.Sleep(time.Duration(deltaTime))
	}
}

func startAgents(t *Target, amount int) {
	for i := 0; i < amount; i++ {
		agent := NewAgent(strconv.Itoa(t.Agents.Len()+1)+"("+t.Name+")", make(chan bool), *factory.CreateGenerator())
		t.Agents.PushBack(agent)
		go agent.Run(t.AgentChannel)
	}
}

func interruptAgents(t *Target, amount int) {
	for i := 0; i < amount; i++ {
		agentElem := t.Agents.Back()
		agent := agentElem.Value.(*Agent)
		t.Agents.Remove(agentElem)
		go agent.Interrupt(t.AgentChannel)
	}
}
