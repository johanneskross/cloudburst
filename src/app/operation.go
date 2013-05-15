package app

import (
	"fmt"
)

const OperationPurchase = 0
const OperationManage = 1

type Operation struct {
	Name                     string
	StartTime, EndTime       int
	Success                  bool
	NumberOfActionsPerformed int
}

func (o Operation) Run() {
	fmt.Printf("%v running..\n", o.Name)
}

type Operate interface {
	Run()
}

func GetOperation(operationId int) Operate {
	switch operationId{
	case OperationPurchase:
		return NewOperationPurchase()
	case OperationManage:
		return NewOperationManage()
	}
	return nil
}

func NewOperationPurchase() Operate {
	op := Operate(Operation{"purchase operation", 0, 0, false, 0})
	return op
}

func NewOperationManage() Operate {
	op := Operate(Operation{"manage operation", 0, 0, false, 0})
	return op
}
