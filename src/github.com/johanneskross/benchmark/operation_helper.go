package benchmark

import (
	"bufio"
	"container/list"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type OperationHelper struct {
	QueryString string
	HttpClient  http.Client
	SourceCode  *list.List
}

func NewOperationHelper(quertyString string) OperationHelper {
	o := OperationHelper{}
	o.QueryString = quertyString
	return o
}

func (o *OperationHelper) Prepare() {

}

func (o *OperationHelper) Cleanup() {

}

func (o *OperationHelper) Login() {
	cookieJar := &CookieJarImpl{}
	cookieJar.cookies = make(map[string][]*http.Cookie)

	o.HttpClient = http.Client{}
	o.HttpClient.Jar = cookieJar

	uid := RandInt(1, dbSize)
	values := url.Values{}
	values.Set("uid", strconv.Itoa(uid))

	resp, _ := o.HttpClient.PostForm(o.QueryString, values)
	resp.Body.Close()
}

func (o *OperationHelper) Logout(queryString, uid string) {

}

func (o *OperationHelper) GoHome() {
	for i := 0; i < RETRY; i++ {
		values := url.Values{}
		values.Set("action", "home")

		resp, err := o.HttpClient.PostForm(o.QueryString, values)
		resp.Body.Close()
		if err != nil {
			break
		}
	}
}

func (o *OperationHelper) BrowseVehicles(direction string, category int) {
	if direction == "top" {
		for i := 0; i < RETRY; i++ {
			values := url.Values{}
			values.Set("action", "View_Items")
			values.Set("category", strconv.Itoa(category))

			resp, err := o.HttpClient.PostForm(o.QueryString, values)
			resp.Body.Close()
			if err != nil {
				break
			}
		}
	} else {
		for i := 0; i < RETRY; i++ {
			values := url.Values{}
			values.Set("action", "View_Items")
			values.Set("browse", direction)
			values.Set("category", strconv.Itoa(category))

			resp, err := o.HttpClient.PostForm(o.QueryString, values)
			resp.Body.Close()
			if err != nil {
				break
			}
		}
	}
}

func (o *OperationHelper) SaveSourceCode(resp *http.Response) {
	lines := list.New()
	reader := bufio.NewReader(resp.Body)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lines.PushBack(string(line))
	}
	o.SourceCode = lines
}

func (o *OperationHelper) DealershipInventory() {
	for i := 0; i < RETRY; i++ {
		values := url.Values{}
		values.Set("action", "inventory")

		resp, err := o.HttpClient.PostForm(o.QueryString, values)
		o.SaveSourceCode(resp)
		resp.Body.Close()
		if err != nil {
			break
		}
	}
}

func (o *OperationHelper) GetOpenOrders() *list.List {
	openOrders := list.New()
	for elem := o.SourceCode.Front(); elem != nil; elem = elem.Next() {
		line := elem.Value.(string)
		if strings.Contains(line, "action=cancelorder") {
			start := strings.Index(line, "cancelorder&orderID=") + 20
			end := strings.Index(line, "cancel order") - 2
			openOrders.PushBack(line[start:end])
		}
	}
	return openOrders
}

type ItemQuantity struct {
	ItemID                  string
	ItemQuantity, ItemTotal int
}

func (o *OperationHelper) GetVehiclesForSale() *list.List {
	items := list.New()
	for elem := o.SourceCode.Front(); elem != nil; elem = elem.Next() {
		line := elem.Value.(string)
		if strings.Contains(line, "action=sellinventory") {

			start := strings.Index(line, "vehicleToSell=") + 14
			end := strings.Index(line, "&total")
			itemID := line[start:end]

			start = strings.Index(line, "quantity ") + 9
			end = strings.Index(line, "-->") - 1
			itemQuantity, _ := strconv.Atoi(line[start:end])

			start = strings.Index(line, "total=") + 6
			end = strings.Index(line, ">Sell")
			itemTotal, _ := strconv.Atoi(line[start:end])

			item := &ItemQuantity{itemID, itemQuantity, itemTotal}
			items.PushBack(item)
		}
	}
	return items
}

func (o *OperationHelper) AddVehiclesToCart(items []ItemQuantity, category int) {
	for i := 0; i < len(items); i++ {

		for j := 0; j < RETRY; j++ {
			values := url.Values{}
			values.Set("action", "View_Items")
			values.Set("vehicles", items[i].ItemID)
			values.Set("quantity", strconv.Itoa(items[i].ItemQuantity))
			values.Set("category", strconv.Itoa(category))

			resp, err := o.HttpClient.PostForm(o.QueryString, values)
			o.SaveSourceCode(resp)
			resp.Body.Close()
			if err != nil {
				break
			}
		}

		details := ""
		for elem := o.SourceCode.Front(); elem != nil; elem = elem.Next() {
			line := elem.Value.(string)
			if strings.Contains(line, "driver-tag-start") {
				details = line
				break
			}
		}

		if details == "" {
			log.Fatal("No driver-tag-start found in response")
		}

		start := strings.Index(details, "name:") + 5
		end := strings.Index(details, "description")
		name := details[start:end]

		start = strings.Index(details, "description:") + 12
		end = strings.Index(details, "price:")
		description := details[start:end]

		start = strings.Index(details, "price:") + 6
		end = strings.Index(details, "discount") - 1
		price := details[start:end]

		start = strings.Index(details, "discount:") + 9
		end = strings.Index(details, "driver-tag-end") - 1
		discount := details[start:end]

		for j := 0; j < RETRY; j++ {
			values := url.Values{}
			values.Set("action", "Add to Cart")
			values.Set("itemId", items[i].ItemID)
			values.Set("quantity", strconv.Itoa(items[i].ItemQuantity))
			values.Set("name", name)
			values.Set("description", description)
			values.Set("price", price)
			values.Set("discount", discount)

			resp, err := o.HttpClient.PostForm(o.QueryString, values)
			o.SaveSourceCode(resp)
			resp.Body.Close()
			if err != nil {
				break
			}
		}

	}
}

func (o *OperationHelper) CancelOpenOrder(orderID string) {
	for i := 0; i < RETRY; i++ {
		values := url.Values{}
		values.Set("action", "cancelorder")
		values.Set("orderID", orderID)

		resp, err := o.HttpClient.PostForm(o.QueryString, values)
		resp.Body.Close()
		if err != nil {
			break
		}

		// TODO check for app error -> concurrency
	}
}

func (o *OperationHelper) SellVehiclesFromInventory(car string, quantaty int, total int) {
	for i := 0; i < RETRY; i++ {
		values := url.Values{}
		values.Set("action", "sellinventory")
		values.Set("vehicleToSell", car)
		values.Set("total", strconv.Itoa(total))

		resp, err := o.HttpClient.PostForm(o.QueryString, values)
		resp.Body.Close()
		if err != nil {
			break
		}
	}
}

func (o *OperationHelper) ClearCart() {
	for i := 0; i < RETRY; i++ {
		values := url.Values{}
		values.Set("action", "shoppingcart")

		resp, err := o.HttpClient.PostForm(o.QueryString, values)
		resp.Body.Close()
		if err != nil {
			break
		}
	}

	for i := 0; i < RETRY; i++ {
		values := url.Values{}
		values.Set("action", "clearcart")

		resp, err := o.HttpClient.PostForm(o.QueryString, values)
		resp.Body.Close()
		if err != nil {
			break
		}
	}
}

func (o *OperationHelper) RemoveVehiclesFromCart(itemID string) {
	for i := 0; i < RETRY; i++ {
		values := url.Values{}
		values.Set("action", "shoppingcart")

		resp, err := o.HttpClient.PostForm(o.QueryString, values)
		resp.Body.Close()
		if err != nil {
			break
		}
	}

	for i := 0; i < RETRY; i++ {
		values := url.Values{}
		values.Set("action", "remove")
		values.Set("cartID", "itemID")

		resp, err := o.HttpClient.PostForm(o.QueryString, values)
		resp.Body.Close()
		if err != nil {
			break
		}
	}
}

func (o *OperationHelper) CheckOut(checktype string, location int) {
	for i := 0; i < RETRY; i++ {
		values := url.Values{}
		if checktype == "purchase" {
			values.Set("action", "purchasecart")
			values.Set("location", strconv.Itoa(location))
		} else {
			values.Set("action", "deferorder")
		}

		resp, err := o.HttpClient.PostForm(o.QueryString, values)
		resp.Body.Close()
		if err != nil {
			break
		}
	}
}
