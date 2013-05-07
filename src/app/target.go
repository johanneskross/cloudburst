package app

import (
	"container/list"
	"fmt"
	"strconv"
	"time"
)

type Target struct {
	Name			string
	Agents			list.List
	ActiveAgents	int
	AgentChannel	chan bool
}

func NewTarget(name string, amountAgents int) *Target {	
	agentChannel := make(chan bool, amountAgents)
	agents := *list.New()
	return &Target{name, agents, amountAgents, agentChannel}
} 

func (t *Target) UpdateAgents() {
	for {
		currentAgents := len(t.AgentChannel)
		switch {
			case currentAgents < t.ActiveAgents:
				newAgents := t.ActiveAgents - currentAgents
				for i := 0; i < newAgents; i++ {
					agent := NewAgent(strconv.Itoa(i), make(chan bool))
					t.Agents.PushBack(agent)
					go agent.Run(t.AgentChannel)
				}
			case currentAgents > t.ActiveAgents:
				reduceAgents := currentAgents - t.ActiveAgents
				for i := 0; i < reduceAgents; i++ {
					agentElem := t.Agents.Back()
					agent := agentElem.Value.(*Agent)
					agent.Quit <- true
					t.Agents.Remove(agentElem)
				}
		}
		time.Sleep(20 * time.Second)
		fmt.Println(len(t.AgentChannel))
	}
}