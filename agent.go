package cloudburst

import (
	//"fmt"
	"github.com/johanneskross/cloudburst/scoreboard"
	"time"
)

type Agent struct {
	AgentId, TargetId      int
	TargetIp               string
	AgentJoinChannel       chan bool
	Generator              Generator
	OperationResultChannel chan *scoreboard.OperationResult
	WaitTimeChannel        chan *scoreboard.WaitTime
	Timing                 *Timing
}

func NewAgent(agentId, targetId int, targetIp string, agentJoinChannel chan bool, generator Generator, operationResultChannel chan *scoreboard.OperationResult, waitTimeChannel chan *scoreboard.WaitTime, timing *Timing) *Agent {
	return &Agent{agentId, targetId, targetIp, agentJoinChannel, generator, operationResultChannel, waitTimeChannel, timing}
}

func (agent *Agent) Run(agentChannel chan bool) {
	//	fmt.Printf("Starting agent #%v ..\n", agent.AgentId)
	for {
		select {
		case <-agent.AgentJoinChannel:
			// fmt.Printf("Stopping agent: #%v ..\n", agent.AgentId)
			close(agent.AgentJoinChannel)
			agentChannel <- true
			return
		default:
			operation := agent.Generator.NextRequest(agent.TargetIp)
			agent.OperateSync(operation)
		}
	}
}

func (agent *Agent) OperateSync(operation Operation) {
	// execute operation
	startTime := time.Now().UnixNano()
	operationResult := operation.Run(agent.Timing)
	endTime := time.Now().UnixNano()

	// report operation result
	operationResult.TimeStarted = startTime
	operationResult.TimeFinished = endTime
	agent.OperationResultChannel <- operationResult

	// wait
	startTime = time.Now().UnixNano()
	waitTime := agent.Generator.GetWaitTime() * TO_NANO
	waitTime = agent.Sync(startTime, waitTime)

	// report wait time
	agent.WaitTimeChannel <- scoreboard.NewWaitTime(startTime, waitTime, operationResult.OperationName)
}

func (agent *Agent) Sync(startTime, waitTime int64) int64 {
	endTime := startTime + waitTime
	if endTime > agent.Timing.EndRun {
		if startTime < agent.Timing.Start {
			waitTime = agent.Timing.SteadyStateDuration() - startTime
			agent.WaitUntil(agent.Timing.StartSteadyState)
		} else {
			waitTime = agent.Timing.EndRun - startTime
			agent.WaitUntil(agent.Timing.EndRun)
		}
	} else {
		agent.WaitUntil(endTime)
	}

	return waitTime
}

func (agent *Agent) WaitUntil(endTime int64) {
	startTime := time.Now().UnixNano()
	duration := endTime - startTime
	if duration > 0 {
		time.Sleep(time.Duration(duration))
	}
}
