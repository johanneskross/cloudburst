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
	return
	o.Do.Login()
	o.Do.DealershipInventory()

	openOrders := o.Do.GetOpenOrders()
	openOrdersCount := openOrders.Len()

	if openOrdersCount >= OPEN_ORDER_DEFAULT_CANCEL_GROUP_SIZE {

		cancelCount := OPEN_ORDER_DEFAULT_CANCEL_GROUP_SIZE
		if cancelCount >= OPEN_ORDER_BIG_CANCEL_GROUP_SIZE {
			cancelCount = openOrdersCount - OPEN_ORDER_BIG_CANCEL_GROUP_SIZE
		}

		e := openOrders.Front()
		for i := 0; i < cancelCount; i++ {
			order := e.Value.(string)
			o.Do.CancelOpenOrder(order)
			e.Next()
		}
	}

	o.Do.DealershipInventory()
	vehicles := o.Do.GetVehiclesForSale()
	if vehicles.Len() != 0 {

		totalVehicles := 0
		for elem := vehicles.Front(); elem != nil; elem = elem.Next() {
			vehicle := elem.Value.(ItemQuantity)
			totalVehicles += vehicle.ItemQuantity
		}

		carsToSell := totalVehicles - VEHICLES_NOT_SOLD
		if carsToSell < 0 {
			carsToSell = 0
		}
		carsSold := 0
		for elem := vehicles.Front(); elem != nil; elem = elem.Next() {
			vehicle := elem.Value.(ItemQuantity)

			o.Do.SellVehiclesFromInventory(vehicle.ItemID, vehicle.ItemQuantity, vehicle.ItemTotal)
			carsSold += vehicle.ItemQuantity
			if carsSold >= carsToSell {
				break
			}
		}
	}
	o.Do.GoHome()
}
