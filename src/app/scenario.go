package app

import (
	"times"
)

type Scenario struct {
}

func createSchedule() *TargetSchedule {
	schedule := NewTargetSchedule()
	conf1 := NewTargetConfiguration("1", 0, 5, 5, 5, *times.LoadTimeSeries())
	conf2 := NewTargetConfiguration("2", 0, 5, 5, 5, *times.LoadTimeSeries())
	conf3 := NewTargetConfiguration("3", 0, 5, 5, 5, *times.LoadTimeSeries())
	schedule.AddTargetConfiguration(conf1)
	schedule.AddTargetConfiguration(conf2)
	schedule.AddTargetConfiguration(conf3)
	return schedule
}

func createManager(schedule TargetSchedule) *TargetManager {
	return NewTargetManager(schedule)
}

func (scenario *Scenario) Launch() {
	schedule := *createSchedule()
	manager := *createManager(schedule)
	joinChannel := make(chan bool)
	go manager.processSchedule(joinChannel)
	<-joinChannel
}
