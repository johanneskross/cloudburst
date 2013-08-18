package cloudburst

import (
	"container/list"
)

type TargetFactory interface {
	CreateTargets(targetConfiguration *TargetConfiguration) *list.List
	CountTargets() int
}
