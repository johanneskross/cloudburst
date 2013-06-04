package benchmark

import (
	"fmt"
)

type OperationPurchase struct {
	Id                       string
	StartTime, EndTime       int
	Success                  bool
	NumberOfActionsPerformed int
	Category                 int
	Do                       OperationHelper
}

func NewOperationPurchase(name string, startTime, endTime int, success bool, numberOfActionsPerformed int, helper OperationHelper) *OperationPurchase {
	return &OperationPurchase{name, startTime, endTime, success, numberOfActionsPerformed, 2, helper}
}

func (o *OperationPurchase) Name() string {
	return o.Id
}

func (o *OperationPurchase) Run() {
	return
	fmt.Printf("running purchase operation..\n")
	orderlineCount := RandInt(MIN_ORDERLINE_COUNT, MAX_ORDERLINE_COUNT)
	items := make([]ItemQuantity, orderlineCount)
	itemIDs := make([]string, orderlineCount)

	totalVehicleQuantity := 0
	isLargeOrder := RandBool(LARGE_ORDER_PERCENTAGE)
	if isLargeOrder {
		totalVehicleQuantity = RandInt(MIN_LARGE_ORDER_VEHICLE_COUNT, MAX_LARGE_ORDER_VEHICLE_COUNT)
	} else {
		totalVehicleQuantity = RandInt(MIN_ORDER_VEHICLE_COUNT, MAX_ORDER_VEHICLE_COUNT)
	}

	isDefferedOrder := !isLargeOrder && RandBool(DEFERRED_ORDER_PERCENTAGE)

	quantity := totalVehicleQuantity / orderlineCount
	rem := totalVehicleQuantity - quantity*orderlineCount
	for i := 0; i < orderlineCount; i++ {
		done := false
		for !done {
			itemIDs[i] = RandPartId()

			j := 0
			for j = 0; j < i; j++ {
				if itemIDs[i] == itemIDs[j] {
					break
				}
			}
			done = j == i
		}
	}

	for i := 0; i < orderlineCount; i++ {
		item := ItemQuantity{itemIDs[i], quantity, 0}
		items[i] = item
	}

	items[orderlineCount-1].ItemQuantity += rem

	o.Do.Login()
	o.Category = RandCategory()
	o.Do.AddVehiclesToCart(items, o.Category)

	if !RandBool(CHECKOUT_ENTIRE_CART_PERCENTAGE) {
		if RandBool(REFILL_CART_PERCENTAGE) {
			o.Do.ClearCart()
			o.Do.AddVehiclesToCart(items, o.Category)
		} else {
			removeOrderlineCount := orderlineCount - REMAINING_ORDERLINE_COUNT
			totalVehicleQuantity = items[removeOrderlineCount].ItemQuantity
			for i := 0; i < removeOrderlineCount; i++ {
				o.Do.RemoveVehiclesFromCart(items[i].ItemID)
			}
		}
	}

	location := 0
	if isLargeOrder {
		randomItemPlace := RandInt(0, numOfItems-1)
		location = randomItemPlace / numItemsPerLoc
	}

	checktype := ""
	if isDefferedOrder {
		checktype = "defer"
	} else {
		checktype = "purchase"
	}
	o.Do.CheckOut(checktype, location)
	o.Do.GoHome()
}
