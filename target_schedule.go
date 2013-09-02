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
	duration := targetConfiguration.Offset + targetConfiguration.RampUp + targetConfiguration.Duration + targetConfiguration.RampDown
	if duration > targetSchedule.Duration {
		targetSchedule.Duration = duration
	}

	if targetSchedule.TargetConfigurations == nil {
		targetSchedule.TargetConfigurations = list.New()
	}
	targetSchedule.TargetConfigurations.PushFront(targetConfiguration)
}
