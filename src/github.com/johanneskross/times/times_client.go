package times

import (
	"github.com/samuel/go-thrift"
	"net"
	"strconv"
)

const NET = "tcp"

func connect(timesHost string, port int) TimeServiceClient {
	address := timesHost + ":" + strconv.Itoa(port)
	conn, err := net.Dial(NET, address)
	if err != nil {
		panic(err)
	}
	client := thrift.NewClient(conn, thrift.NewBinaryProtocol(true, false))

	return TimeServiceClient{client}
}

func LoadTimeSeries(timesHost string, port int, timeSeriesName string) *TimeSeries {
	time := connect(timesHost, port)
	timeSeries, err := time.Load(timeSeriesName)
	if err != nil {
		panic(err)
	}
	return timeSeries
}
