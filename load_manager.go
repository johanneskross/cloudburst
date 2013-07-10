package cloudburst

import ()

type LoadManager struct {
	TargetId        int
	LoadSchedule    *LoadSchedule
	CurrentLoadUnit *LoadUnit
	NextLoadIndex   int
}

func NewLoadManager(loadSchedule *LoadSchedule) *LoadManager {
	return &LoadManager{-1, loadSchedule, nil, -1}
}

func (loadManager *LoadManager) NextLoadUnit() *LoadUnit {
	loadManager.NextLoadIndex = loadManager.NextLoadIndex + 1
	nextLoadUnit := loadManager.LoadSchedule.GetLoadUnit(loadManager.NextLoadIndex)
	return nextLoadUnit
}
