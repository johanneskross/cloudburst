package cloudburst

import ()

type LoadManager struct {
	TargetId        int
	LoadSchedule    *LoadSchedule
	CurrentLoadUnit *LoadUnit
	NextLoadIndex   int
}

func NewLoadManager(targetId int, loadSchedule *LoadSchedule) *LoadManager {
	return &LoadManager{targetId, loadSchedule, nil, -1}
}

func (loadManager *LoadManager) NextLoadUnit() *LoadUnit {
	loadManager.NextLoadIndex = loadManager.NextLoadIndex + 1
	nextLoadUnit := loadManager.LoadSchedule.GetLoadUnit(loadManager.NextLoadIndex)
	return nextLoadUnit
}
