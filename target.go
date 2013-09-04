package cloudburst

import (
	"container/list"
	"fmt"
	"github.com/johanneskross/cloudburst/load"
	"time"
)

const TO_NANO = 1000000000

type Target struct {
	TargetId              int
	Agents                *list.List
	AgentChannel          chan bool
	Configuration         *TargetConfiguration
	Generator             Generator
	Scoreboard            *Scoreboard
	ScoreboardJoinChannel chan bool
	Timing                *Timing
	LoadManager           *load.LoadManager
}

func NewTarget(targetConfiguration *TargetConfiguration, generator Generator, loadManager *load.LoadManager) *Target {
	target := &Target{}
	target.TargetId = -1
	target.Agents = list.New()
	target.AgentChannel = make(chan bool, loadManager.LoadSchedule.MaxAgents())
	target.Configuration = targetConfiguration
	target.Generator = generator
	target.LoadManager = loadManager
	return target
}

func (t *Target) GenerateLoad(targetJoinChannel chan bool) {
	t.createAndStartScoreboard()

	// wait for start
	t.WaitUntil(t.Timing.StartSteadyState)

	fmt.Printf("Running time series on target %v with ip %v\n",
		t.TargetId, t.Configuration.TargetIp)

	interval := 0
	for t.Timing.InSteadyState(time.Now().UnixNano()) {
		// get load unit
		loadUnit := t.LoadManager.NextLoadUnit()
		loadUnit.Activiate()

		// update amount of agents
		t.updateAmountOfAgents(loadUnit, interval)

		// wait until interval ends
		t.WaitUntil(loadUnit.IntervalEnd())
		interval++
	}
	// terminate target
	t.terminateTarget(targetJoinChannel)
}

func (t *Target) updateAmountOfAgents(loadUnit *load.LoadUnit, interval int) {
	runningAgents := t.Agents.Len()
	runningNextAgents := int(loadUnit.NumberOfUsers)
	runningNextAgents = 150

	fmt.Printf("Update amount of agents from %v to %v on target%v in interval %v\n",
		runningAgents, runningNextAgents, t.TargetId, interval)

	// update amount of agents for this interval
	switch {
	case runningAgents < runningNextAgents:
		addAgents := runningNextAgents - runningAgents
		t.startAgents(addAgents)
	case runningAgents > runningNextAgents:
		reduceAgents := runningAgents - runningNextAgents
		t.interruptAgents(reduceAgents, false)
	}
}

func (t *Target) WaitUntil(nextInterval int64) {
	currentTime := time.Now().UnixNano()
	deltaTime := nextInterval - currentTime

	if deltaTime > 0 {
		fmt.Printf("Target %v waits %v seconds for next interval\n",
			t.TargetId, deltaTime/TO_NANO)
		time.Sleep(time.Duration(deltaTime))
	}
}

func (t *Target) startAgents(amount int) {
	for i := 0; i < amount; i++ {
		agent := &Agent{}
		agent.AgentId = t.Agents.Len() + 1
		agent.TargetId = t.TargetId
		agent.TargetIp = t.Configuration.TargetIp
		agent.AgentJoinChannel = make(chan bool, 1)
		agent.Generator = t.Generator
		agent.OperationResultChannel = t.Scoreboard.OperationResultChannel
		agent.WaitTimeChannel = t.Scoreboard.WaitTimeChannel
		agent.Timing = t.Timing

		t.Agents.PushBack(agent)
		go agent.Run(t.AgentChannel)
	}
}

func (t *Target) interruptAgents(amount int, waitForInterrupt bool) {
	for i := 0; i < amount; i++ {
		agentElem := t.Agents.Back()
		agent := agentElem.Value.(*Agent)
		agent.AgentJoinChannel <- true
		t.Agents.Remove(agentElem)
	}
	if waitForInterrupt {
		for i := 0; i < amount; i++ {
			<-t.AgentChannel
		}
	}
}

func (t *Target) createAndStartScoreboard() {
	t.Scoreboard = NewScoreboard(t.TargetId, cap(t.AgentChannel), t.Timing)
	t.ScoreboardJoinChannel = make(chan bool)
	go t.Scoreboard.Run(t.ScoreboardJoinChannel)
}

func (t *Target) terminateTarget(targetJoinChannel chan bool) {
	t.interruptAgents(t.Agents.Len(), true)
	t.ScoreboardJoinChannel <- true
	<-t.ScoreboardJoinChannel
	targetJoinChannel <- true
}
