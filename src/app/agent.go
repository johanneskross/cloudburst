package app

import (
	"fmt"
	"time"
)

const URL = "http://www.example.com/"

type Agent struct {
	Name, Url string
	Quit      chan bool
}

func NewAgent(name string, quit chan bool) *Agent {
	return &Agent{name, URL, quit}
}

func (agent *Agent) GenerateLoad() {
	pullUrl(agent.Url)
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
			agent.GenerateLoad()
			time.Sleep(1 * time.Second)
		}
	}
}

func (agent *Agent) Interrupt(c chan bool) {
	agent.Quit <- true
}
