package benchmark

import (
	"fmt"
)

type OperationManage struct {
	Id                       string
	StartTime, EndTime       int
	Success                  bool
	NumberOfActionsPerformed int
	Do                       OperationHelper
}

func NewOperationManage(name string, startTime, endTime int, success bool, numberOfActionsPerformed int, helper OperationHelper) *OperationManage {
	return &OperationManage{name, startTime, endTime, success, numberOfActionsPerformed, helper}
}

func (o *OperationManage) Name() string {
	return o.Id
}

func (o *OperationManage) Run() {
	fmt.Printf("running manage operation..\n")
	o.Do.Prepare()
}
