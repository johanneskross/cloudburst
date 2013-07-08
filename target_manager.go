package cloudburst

import ()

type TargetManager struct {
	Schedule TargetSchedule
	Factory  Factory
	Targets  []*Target
}

func NewTargetManager(schedule TargetSchedule, factory Factory) *TargetManager {
	return &TargetManager{schedule, factory, nil}
}

func (targetManager *TargetManager) processSchedule(joinChannel chan bool) {
	targetConfigurations := targetManager.Schedule.TargetConfigurations
	joinTargetChannel := make(chan bool, targetConfigurations.Len())
	targetManager.Targets = make([]*Target, targetConfigurations.Len())

	// start targets
	i := 0
	for elem := targetConfigurations.Front(); elem != nil; elem = elem.Next() {
		targetConfiguration := elem.Value.(*TargetConfiguration)
		target := NewTarget(*targetConfiguration, targetManager.Factory)
		targetManager.Targets[i] = target
		go target.RunTimeSeries(joinTargetChannel)
		i++
	}

	// wait until all targets ended
	for i := 0; i < cap(joinTargetChannel); i++ {
		<-joinTargetChannel
	}

	joinChannel <- true
}
