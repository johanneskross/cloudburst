package benchmark

import (
	"fmt"
	"math"
)

type OperationBrowse struct {
	Id                       string
	StartTime, EndTime       int
	Success                  bool
	NumberOfActionsPerformed int
	Category                 int
	Do                       OperationHelper
}

func NewOperationBrowse(name string, startTime, endTime int, success bool, numberOfActionsPerformed int, helper OperationHelper) *OperationBrowse {
	return &OperationBrowse{name, startTime, endTime, success, numberOfActionsPerformed, 2, helper}
}

func (o *OperationBrowse) Name() string {
	return o.Id
}

func (o *OperationBrowse) Run() {
	fmt.Printf("running browse operation..\n")
	return

	o.Do.Login()
	o.Category = RandCategory()
	o.Do.BrowseVehicles("top", o.Category)

	for i := 0; i < FORWARD_BROWSES; i++ {
		o.Do.BrowseVehicles("fwd", o.Category)

		if math.Mod(float64(i), BACKWARD_BROWSE_INTERVAL) == BACKWARD_BROWSE_INTERVAL-1 {
			o.Do.BrowseVehicles("bwd", o.Category)
		}
	}

	o.Do.GoHome()
}
