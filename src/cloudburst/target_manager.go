package cloudburst

import (
	"fmt"
)

type TargetManager struct {
	Schedule TargetSchedule
}

func NewTargetManager(schedule TargetSchedule) *TargetManager {
	return &TargetManager{schedule}
}

func (targetManager *TargetManager) processSchedule(joinChannel chan bool) {
	fmt.Printf("Starting target manager\n")

	targetConfigurations := targetManager.Schedule.TargetConfigurations
	joinTargetChannel := make(chan bool, targetConfigurations.Len())

	// start targets
	for elem := targetConfigurations.Front(); elem != nil; elem = elem.Next() {
		targetConfiguration := elem.Value.(*TargetConfiguration)
		target := NewTarget(*targetConfiguration)
		go target.RunTimeSeries(joinTargetChannel)
	}
	
	// wait until all targets ended
	for i := 0; i < cap(joinTargetChannel); i++ {
		<-joinTargetChannel
	}

	joinChannel <- true
	fmt.Printf("Ending target manager\n")
}
