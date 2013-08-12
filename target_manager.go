package cloudburst

import (
	"container/list"
	"time"
)

type TargetManager struct {
	Schedule *TargetSchedule
	Factory  Factory
	Targets  *list.List
	TargetId int
}

func NewTargetManager(schedule *TargetSchedule, factory Factory) *TargetManager {
	return &TargetManager{schedule, factory, list.New(), 0}
}

func (targetManager *TargetManager) processSchedule(joinChannel chan bool) {
	targetConfigurations := targetManager.Schedule.TargetConfigurations

	joinTargetChannel := make(chan bool, targetManager.countAllTargets())

	startBenchmarkTime := time.Now().UnixNano()
	// create and start targets
	for elem := targetConfigurations.Front(); elem != nil; elem = elem.Next() {
		targetConfiguration := elem.Value.(*TargetConfiguration)
		targetManager.createAndStartTarget(targetConfiguration, joinTargetChannel, startBenchmarkTime)
	}

	// wait until all targets ended
	for i := 0; i < cap(joinTargetChannel); i++ {
		<-joinTargetChannel
	}

	joinChannel <- true
}

func (targetManager *TargetManager) createAndStartTarget(targetConfiguration *TargetConfiguration, joinTargetChannel chan bool, startBenchmarkTime int64) {
	waitTime := startBenchmarkTime + targetConfiguration.Offset - time.Now().UnixNano()
	time.Sleep(time.Duration(waitTime))

	targets := targetConfiguration.TargetFactory.CreateTargets(targetConfiguration, targetManager.Factory)
	for elem := targets.Front(); elem != nil; elem = elem.Next() {
		target := elem.Value.(*Target)

		// Set target values
		target.TargetId = targetManager.TargetId
		targetManager.TargetId++
		target.LoadManager.TargetId = targetManager.TargetId
		target.Timing = NewTiming(targetConfiguration.RampUp, targetConfiguration.Duration, targetConfiguration.RampDown)

		targetManager.Targets.PushBack(target)
		go target.RunTimeSeries(joinTargetChannel)
	}
}

func (targetManager *TargetManager) countAllTargets() int {
	targetConfigurations := targetManager.Schedule.TargetConfigurations
	targetCount := 0
	for elem := targetConfigurations.Front(); elem != nil; elem = elem.Next() {
		targetConfiguration := elem.Value.(*TargetConfiguration)
		targetCount += targetConfiguration.TargetFactory.CountTargets()
	}
	return targetCount
}
