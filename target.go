package cloudburst

import (
	"container/list"
	"fmt"
	"github.com/johanneskross/cloudburst/load"
	"time"
)

const TO_NANO = 1000000000

type Target struct {
	TargetId      int
	Agents        *list.List
	AgentChannel  chan bool
	Configuration *TargetConfiguration
	Generator     Generator
	Scoreboard    *Scoreboard
	Timing        *Timing
	LoadManager   *load.LoadManager
}

func NewTarget(targetConfiguration *TargetConfiguration, generator Generator, loadManager *load.LoadManager) *Target {
	channelSize := loadManager.LoadSchedule.MaxAgents()
	channelSize = 200 // For test purpose
	agentChannel := make(chan bool, channelSize)

	target := &Target{}
	target.TargetId = -1
	target.Agents = list.New()
	target.AgentChannel = agentChannel
	target.Configuration = targetConfiguration
	target.Generator = generator
	target.LoadManager = loadManager
	return target
}

func (t *Target) RunTimeSeries(c chan bool) {
	fmt.Printf("Running time series on target %v with ip %v\n", t.TargetId, t.Configuration.TargetIp)

	t.Scoreboard = NewScoreboard(t.TargetId, t.Timing)
	scoreboardQuitQuannel := make(chan bool)
	go t.Scoreboard.Run(scoreboardQuitQuannel)

	t.WaitUntil(t.Timing.StartSteadyState)

	i := 0
	for t.Timing.InSteadyState(time.Now().UnixNano()) {
		// wait until next interval is due
		loadUnit := t.LoadManager.NextLoadUnit()
		loadUnit.Activiate()

		runningAgents := t.Agents.Len()
		runningNextAgents := int(loadUnit.NumberOfUsers)
		//runningNextAgents = 200 // For test purpose
		fmt.Printf("Update amount of agents from %v to %v on target%v in interval %v\n", runningAgents, runningNextAgents, t.TargetId, i)

		// update amount of agents for this interval
		switch {
		case runningAgents < runningNextAgents:
			addAgents := runningNextAgents - runningAgents
			startAgents(t, addAgents)
		case runningAgents > runningNextAgents:
			reduceAgents := runningAgents - runningNextAgents
			go interruptAgents(t, reduceAgents)
		}

		t.WaitUntil(loadUnit.IntervalEnd())
		i++
	}
	interruptAgents(t, t.Agents.Len())
	scoreboardQuitQuannel <- true
	<-scoreboardQuitQuannel
	c <- true
}

func (t *Target) WaitUntil(nextInterval int64) {
	currentTime := time.Now().UnixNano()
	deltaTime := nextInterval - currentTime

	if deltaTime > 0 {
		fmt.Printf("Target %v waits %v seconds for next interval\n", t.TargetId, deltaTime/TO_NANO)
		time.Sleep(time.Duration(deltaTime))
	}
}

func startAgents(t *Target, amount int) {
	for i := 0; i < amount; i++ {
		agent := NewAgent(t.Agents.Len()+1, t.TargetId, t.Configuration.TargetIp, make(chan bool, 1), t.Generator, t.Scoreboard.OperationResultChannel, t.Scoreboard.WaitTimeChannel, t.Timing)
		t.Agents.PushBack(agent)
		go agent.Run(t.AgentChannel)
	}
}

func interruptAgents(t *Target, amount int) {
	for i := 0; i < amount; i++ {
		agentElem := t.Agents.Back()
		agent := agentElem.Value.(*Agent)
		agent.Quit <- true
		t.Agents.Remove(agentElem)
	}
	for i := 0; i < amount; i++ {
		<-t.AgentChannel
	}
}
