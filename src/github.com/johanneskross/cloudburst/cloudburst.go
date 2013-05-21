package cloudburst

import (
)

type Scenario struct {
}

func NewScenario() *Scenario{
	return &Scenario{}
}

func createManager(schedule TargetSchedule, factory Factory) *TargetManager {
	return NewTargetManager(schedule, factory)
}

func (scenario *Scenario) Launch(schedule TargetSchedule, factory Factory) {
	manager := *createManager(schedule, factory)
	joinChannel := make(chan bool)
	go manager.processSchedule(joinChannel)
	<-joinChannel
}
