package cloudburst

import ()

type Scenario struct {
	TargetManager *TargetManager
}

func NewScenario(schedule TargetSchedule, factory Factory) *Scenario {
	return &Scenario{NewTargetManager(schedule, factory)}
}

func (scenario *Scenario) Launch() {
	joinChannel := make(chan bool)
	go scenario.TargetManager.processSchedule(joinChannel)
	<-joinChannel
}

func (scenario *Scenario) AggregateStatistics() {
	AggregateScoreboards(scenario.TargetManager.Targets, 0)
}
