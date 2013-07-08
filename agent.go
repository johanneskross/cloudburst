package cloudburst

import (
	//"fmt"
	"time"
)

type Agent struct {
	AgentId, TargetId      int
	TargetIp               string
	Quit                   chan bool
	Generator              Generator
	OperationResultChannel chan OperationResult
	Timing                 *Timing
}

func NewAgent(agentId, targetId int, targetIp string, quit chan bool, generator Generator, operationResultChannel chan OperationResult, timing *Timing) *Agent {
	return &Agent{agentId, targetId, targetIp, quit, generator, operationResultChannel, timing}
}

func (agent *Agent) Run(c chan bool) {
	//fmt.Printf("Starting agent #%v ..\n", agent.AgentId)

	c <- true

	for {
		select {
		case <-agent.Quit:
			//fmt.Printf("Stopping agent: #%v ..\n", agent.AgentId)
			close(agent.Quit)
			<-c
			return
		default:
			operation := agent.Generator.NextRequest(agent.TargetIp)
			agent.OperateSync(operation)
		}
	}
}

func (agent *Agent) Interrupt(c chan bool) {
	agent.Quit <- true
}

func (agent *Agent) OperateSync(operation Operation) {
	timeStarted := time.Duration(time.Now().UnixNano())

	operationResult := operation.Run(agent.Timing)

	timeFinished := time.Duration(time.Now().UnixNano())

	operationResult.TimeStarted = int64(timeStarted)
	operationResult.TimeFinished = int64(timeFinished)
	agent.OperationResultChannel <- operationResult

	agent.Sync()
}

func (agent *Agent) Sync() {
	thinkTime := time.Duration(agent.Generator.GetThinkTime()) * time.Second
	now := time.Duration(time.Now().Unix()) // in seconds
	endTime := now + thinkTime
	deltaTime := endTime - now
	if deltaTime > 0 {
		time.Sleep(deltaTime)
	}
}
