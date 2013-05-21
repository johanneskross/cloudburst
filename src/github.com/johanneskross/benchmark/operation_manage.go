package benchmark 

import (
	"fmt"
)

type OperationManage struct {
	Id                     string
	StartTime, EndTime       int
	Success                  bool
	NumberOfActionsPerformed int
}

func NewOperationManage(name string, startTime, endTime int, success bool, numberOfActionsPerformed int) *OperationManage{
	return &OperationManage{name, startTime, endTime, success, numberOfActionsPerformed}
}

func (o *OperationManage) Name() string{
	return o.Id
}

func (o *OperationManage) Run() {
	fmt.Printf("running manage operation..\n")
}