package app

import (
	"fmt"
	"time"
)

type Target struct {
	Name			string
	ActiveAgents	int
	AgentChannel	chan bool
}

func NewTarget(name string, amountAgents int) *Target {	
	c := make(chan bool, amountAgents)
	
	return &Target{name, amountAgents, c}
} 

func (t *Target) UpdateAgents() {
	if t.ActiveAgents < 0 {
		// error
	}
	currentAgents := len(t.AgentChannel)
	for {
		switch {
			case currentAgents < t.ActiveAgents:
				newAgents := t.ActiveAgents - currentAgents
				for i := 0; i < newAgents; i++ {
					agent := NewAgent("a")
					go agent.Run(t.AgentChannel, make(chan bool))
				}
			case currentAgents > t.ActiveAgents:
				// do something
			case currentAgents == t.ActiveAgents:
				// do something
		}
		time.Sleep(20 * time.Second)
		fmt.Println(len(t.AgentChannel))
	}
}