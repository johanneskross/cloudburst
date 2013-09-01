package cloudburst

import (
	"container/list"
	"time"
)

type TargetManager struct {
	Schedule *TargetSchedule
	Targets  *list.List
	TargetId int
}

func NewTargetManager(schedule *TargetSchedule) *TargetManager {
	return &TargetManager{schedule, list.New(), 0}
}

func (targetManager *TargetManager) processSchedule(targetManagerJoinChannel chan bool) {
	targetConfigurations := targetManager.Schedule.TargetConfigurations
	targetJoinChannel := make(chan bool, targetManager.countAllTargets())
	startTime := time.Now().UnixNano()

	// create and start targets
	for elem := targetConfigurations.Front(); elem != nil; elem = elem.Next() {
		targetConfiguration := elem.Value.(*TargetConfiguration)
		targetManager.createAndStartTarget(targetConfiguration, targetJoinChannel, startTime)
	}

	// wait until all targets ended
	for i := 0; i < cap(targetJoinChannel); i++ {
		<-targetJoinChannel
	}

	// signal termination
	targetManagerJoinChannel <- true
}

func (targetManager *TargetManager) createAndStartTarget(targetConfiguration *TargetConfiguration, targetJoinChannel chan bool, startTime int64) {
	waitTime := startTime + targetConfiguration.Offset - time.Now().UnixNano()
	time.Sleep(time.Duration(waitTime))

	targets := targetConfiguration.TargetFactory.CreateTargets(targetConfiguration)
	for elem := targets.Front(); elem != nil; elem = elem.Next() {
		target := elem.Value.(*Target)
		targetManager.updateTargetVariables(target, targetConfiguration)
		targetManager.Targets.PushBack(target)
		go target.GenerateLoad(targetJoinChannel)
	}
}

func (targetManager *TargetManager) updateTargetVariables(target *Target, targetConfiguration *TargetConfiguration) {
	target.TargetId = targetManager.TargetId
	targetManager.TargetId++
	target.LoadManager.TargetId = targetManager.TargetId
	target.Timing = NewTiming(targetConfiguration.RampUp, targetConfiguration.Duration, targetConfiguration.RampDown)
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
