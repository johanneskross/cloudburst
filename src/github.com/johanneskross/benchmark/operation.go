package benchmark

import (
	"github.com/johanneskross/cloudburst"
)

func GetOperation(operationId int) cloudburst.Operation {
	switch operationId {
	case OperationBrowseId:
		return CreateOperationBrowse()
	case OperationManageId:
		return CreateOperationManage()
	case OperationPurchaseId:
		return CreateOperationPurchase()
	}
	return CreateOperationBrowse()
}

func CreateOperationBrowse() cloudburst.Operation {
	op := NewOperationBrowse("purchase operation", 0, 0, false, 0, NewOperationHelper("http://localhost:3306/specj/app?"))
	return cloudburst.Operation(op)
}

func CreateOperationManage() cloudburst.Operation {
	op := NewOperationManage("manage operation", 0, 0, false, 0, NewOperationHelper("http://localhost:3306/specj/app?"))
	return cloudburst.Operation(op)
}

func CreateOperationPurchase() cloudburst.Operation {
	op := NewOperationPurchase("manage operation", 0, 0, false, 0, NewOperationHelper("http://localhost:3306/specj/app?"))
	return cloudburst.Operation(op)
}
