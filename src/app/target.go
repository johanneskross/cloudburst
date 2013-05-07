package app

import (
	"./times"
	"container/list"
	"fmt"
	"strconv"
	"time"
)

type Target struct {
	Name         string
	Agents       list.List
	AgentChannel chan bool
	TS           times.TimeSeries
}

func NewTarget(name string, timeSeries times.TimeSeries) *Target {
	agents := *list.New()
	channelSize := calcChannelSize(timeSeries.Elements)
	agentChannel := make(chan bool, channelSize)
	return &Target{name, agents, agentChannel, timeSeries}
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

func (t *Target) RunTimeSeries() {
	duration := len(t.TS.Elements)
	for i := 0; i < duration; i++ {
		runningAgents := len(t.AgentChannel)
		runningNextAgents := int(t.TS.Elements[i].Value)
		fmt.Printf("Update amount of agents to %v\n", runningNextAgents)
		switch {
		case runningAgents < runningNextAgents:
			addAgents := runningNextAgents - runningAgents
			startAgents(t, addAgents)
		case runningAgents > runningNextAgents:
			reduceAgents := runningAgents - runningNextAgents
			interruptAgents(t, reduceAgents)
		}
		freq := t.TS.Frequency
		time.Sleep(time.Duration(freq) * time.Second)
	}
}

func startAgents(t *Target, amount int) {
	for i := 0; i < amount; i++ {
		agent := NewAgent(strconv.Itoa(t.Agents.Len()+1), make(chan bool))
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
