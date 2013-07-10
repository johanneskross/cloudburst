package cloudburst

import (
	"container/list"
)

type TargetSchedule struct {
	TargetConfigurations *list.List
	Duration             int64
}

func NewTargetSchedule() *TargetSchedule {
	return &TargetSchedule{}
}

func (targetSchedule *TargetSchedule) AddTargetConfiguration(targetConfiguration *TargetConfiguration) {
	finishTime := targetConfiguration.Offset + targetConfiguration.RampUp + targetConfiguration.Duration + targetConfiguration.RampDown
	if finishTime > targetSchedule.Duration {
		targetSchedule.Duration = finishTime
	}

	if targetSchedule.TargetConfigurations == nil {
		targetSchedule.TargetConfigurations = list.New()
	}
	targetSchedule.TargetConfigurations.PushFront(targetConfiguration)
}
