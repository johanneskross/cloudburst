package cloudburst

import ()

type Operation interface {
	Name() string
//	StartTime()
//	EndTime()
//	Success() bool
//	NumberOfActionsPerformed() int
	Run()
}
