package cloudburst

import (
	"container/list"
)

type LoadSchedule struct {
	LoadUnits *list.List
}

func (loadSchedule *LoadSchedule) GetLoadUnit(index int) *LoadUnit {
	i := 0
	for elem := loadSchedule.LoadUnits.Front(); elem != nil; elem = elem.Next() {
		if i == index {
			return elem.Value.(*LoadUnit)
		}
		i++
	}
	return loadSchedule.LoadUnits.Front().Value.(*LoadUnit)
}

func (loadSchedule *LoadSchedule) Size() int {
	return loadSchedule.LoadUnits.Len()
}

func (loadSchedule *LoadSchedule) MaxAgents() int {
	loadUnits := loadSchedule.LoadUnits
	max := int64(0)
	for elem := loadUnits.Front(); elem != nil; elem = elem.Next() {
		loadUnit := elem.Value.(*LoadUnit)
		if loadUnit.NumberOfUsers > max {
			max = loadUnit.NumberOfUsers
		}
	}
	return int(max)
}

func NewLoadSchedule(loadUnits *list.List) *LoadSchedule {
	return &LoadSchedule{loadUnits}
}
