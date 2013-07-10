package cloudburst

import (
	"container/list"
)

type TargetFactory interface {
	Configure()
	CreateTargets(targetConfiguration *TargetConfiguration, factory Factory) *list.List
	CountTargets() int
}
