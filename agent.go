package cloudburst

import (
	//	"fmt"
	"time"
)

type Agent struct {
	AgentId, TargetId      int
	TargetIp               string
	Quit                   chan bool
	Generator              Generator
	OperationResultChannel chan *OperationResult
	WaitTimeChannel        chan *WaitTime
	Timing                 *Timing
}

func NewAgent(agentId, targetId int, targetIp string, quit chan bool, generator Generator, operationResultChannel chan *OperationResult, waitTimeChannel chan *WaitTime, timing *Timing) *Agent {
	return &Agent{agentId, targetId, targetIp, quit, generator, operationResultChannel, waitTimeChannel, timing}
}

func (agent *Agent) Run(agentChannel chan bool) {
	//	fmt.Printf("Starting agent #%v ..\n", agent.AgentId)
	for {
		select {
		case <-agent.Quit:
			//			fmt.Printf("Stopping agent: #%v ..\n", agent.AgentId)
			close(agent.Quit)
			agentChannel <- true
			return
		default:
			operation := agent.Generator.NextRequest(agent.TargetIp)
			agent.OperateSync(operation)
		}
	}
}

func (agent *Agent) OperateSync(operation Operation) {
	timeStarted := time.Duration(time.Now().UnixNano())

	operationResult := operation.Run(agent.Timing)

	timeFinished := time.Duration(time.Now().UnixNano())

	operationResult.TimeStarted = int64(timeStarted)
	operationResult.TimeFinished = int64(timeFinished)
	agent.OperationResultChannel <- operationResult

	now := time.Now().UnixNano()
	thinkTime := agent.Generator.GetThinkTime() * TO_NANO
	waitTime := agent.Sync(now, thinkTime)
	agent.WaitTimeChannel <- NewWaitTime(now, waitTime, operationResult.OperationName)
}

func (agent *Agent) Sync(now, thinkTime int64) int64 {
	waitTime := thinkTime
	endTime := now + waitTime

	if endTime > agent.Timing.EndRun {

		if now < agent.Timing.Start {
			waitTime = agent.Timing.SteadyStateDuration() - now
			agent.SleepUntil(agent.Timing.StartSteadyState)
		} else {
			waitTime = agent.Timing.EndRun - now
			agent.SleepUntil(agent.Timing.EndRun)
		}

	} else {
		agent.SleepUntil(waitTime)
	}

	return waitTime
}

func (agent *Agent) SleepUntil(endTime int64) {
	now := time.Now().UnixNano()
	duration := endTime - now
	if duration > 0 {
		time.Sleep(time.Duration(duration))
	}
}
