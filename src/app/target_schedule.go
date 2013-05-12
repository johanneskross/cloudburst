package app

import (
	"container/list"
)

type TargetSchedule struct {
	TargetConfigurations list.List
}

func NewTargetSchedule() *TargetSchedule {
	targetConfigurations := *list.New()
	return &TargetSchedule{targetConfigurations}
}

func (targetSchedule *TargetSchedule) AddTargetConfiguration(targetConfiguration *TargetConfiguration) {
	targetSchedule.TargetConfigurations.PushBack(targetConfiguration)
}
