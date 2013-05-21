package benchmark

import (
	"github.com/johanneskross/cloudburst"
)

const OperationBrowseId = 0
const OperationManageId = 1

func GetOperation(operationId int) cloudburst.Operation {
	switch operationId{
	case OperationBrowseId:
		return CreateOperationBrowse()
	case OperationManageId:
		return CreateOperationManage()
	}
	return CreateOperationBrowse()
}

func CreateOperationBrowse() cloudburst.Operation {
	op := NewOperationBrowse("purchase operation", 0, 0, false, 0)
	return cloudburst.Operation(op)
}

func CreateOperationManage() cloudburst.Operation {
	op := NewOperationManage("manage operation", 0, 0, false, 0)
	return cloudburst.Operation(op)
}
