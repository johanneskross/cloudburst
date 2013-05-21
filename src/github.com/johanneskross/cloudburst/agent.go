package cloudburst

import (
	"fmt"
	"time"
)

const URL = "http://www.example.com/"

type Agent struct {
	Name, Url string
	Quit      chan bool
	Generator Generator
}

func NewAgent(name string, quit chan bool, generator Generator) *Agent {
	return &Agent{name, URL, quit, generator}
}

func (agent *Agent) Run(c chan bool) {
	fmt.Printf("Starting agent #%v ..\n", agent.Name)

	c <- true

	for {
		select {
		case <-agent.Quit:
			fmt.Printf("Stopping agent: #%v ..\n", agent.Name)
			close(agent.Quit)
			<-c
			return
		default:
			operation := agent.Generator.NextRequest()
			agent.OperateSync(operation)
		}
	}
}

func (agent *Agent) Interrupt(c chan bool) {
	agent.Quit <- true
}

func (agent *Agent) OperateSync(operation Operation) {
	operation.Run()
	agent.Sync()
}

func (agent *Agent) Sync() {
	thinkTime := time.Duration(agent.Generator.GetThinkTime()) * time.Second
	now := time.Duration(time.Now().Unix())
	endTime := now + thinkTime
	deltaTime := endTime - now
	if deltaTime > 0 {
		time.Sleep(deltaTime)
	}
}
