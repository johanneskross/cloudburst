package cloudburst

import (
	"container/list"
	"fmt"
	"time"
)

const TO_NANO = 1000000000

type Target struct {
	TargetId      int
	Agents        list.List
	AgentChannel  chan bool
	Configuration *TargetConfiguration
	Factory       Factory
	Scoreboard    *Scoreboard
	Timing        *Timing
	LoadManager   *LoadManager
}

func NewTarget(targetConfiguration *TargetConfiguration, factory Factory, loadManager *LoadManager) *Target {
	channelSize := loadManager.LoadSchedule.MaxAgents()
	channelSize = 200 // For test purpose
	agentChannel := make(chan bool, channelSize)

	target := &Target{}
	target.TargetId = -1
	target.Agents = *list.New()
	target.AgentChannel = agentChannel
	target.Configuration = targetConfiguration
	target.Factory = factory
	target.LoadManager = loadManager
	return target
}

func (t *Target) RunTimeSeries(c chan bool) {
	fmt.Printf("Running time series on target: %v\n", t.TargetId)

	t.Scoreboard = NewScoreboard(t.TargetId)
	scoreboardQuitQuannel := make(chan bool)
	go t.Scoreboard.Run(scoreboardQuitQuannel)

	t.Wait(t.Timing.StartSteadyState)

	for t.Timing.InSteadyState(time.Now().UnixNano()) {
		// wait until next interval is due
		loadUnit := t.LoadManager.NextLoadUnit()
		loadUnit.Activiate()

		runningAgents := len(t.AgentChannel)
		runningNextAgents := int(loadUnit.NumberOfUsers)
		//runningNextAgents = 50 // For test purpose
		fmt.Printf("Update amount of agents to %v on target%v\n", runningNextAgents, t.TargetId)

		// update amount of agents for this interval
		switch {
		case runningAgents < runningNextAgents:
			addAgents := runningNextAgents - runningAgents
			startAgents(t, addAgents)
		case runningAgents > runningNextAgents:
			reduceAgents := runningAgents - runningNextAgents
			interruptAgents(t, reduceAgents)
		}

		t.Wait(loadUnit.IntervalEnd())
	}
	scoreboardQuitQuannel <- true
	<-scoreboardQuitQuannel
	c <- true
}

func (t *Target) Wait(nextInterval int64) {
	currentTime := time.Now().UnixNano()
	deltaTime := nextInterval - currentTime

	if deltaTime > 0 {
		fmt.Printf("Target %v waits %v seconds for next interval\n", t.TargetId, deltaTime/TO_NANO)
		time.Sleep(time.Duration(deltaTime))
	}
}

func startAgents(t *Target, amount int) {
	for i := 0; i < amount; i++ {
		agent := NewAgent(t.Agents.Len()+1, t.TargetId, t.Configuration.TargetIp, make(chan bool), t.Factory.CreateGenerator(), t.Scoreboard.OperationResultChannel, t.Timing)
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
