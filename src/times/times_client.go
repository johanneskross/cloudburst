package times

import (
	"github.com/samuel/go-thrift"
	"net"
)

const NET = "tcp"
const ADDRESS = "monitor0.dfg:7855"
const TS_NAME = "O2_business_ADDORDER_profile_user"

func connect() TimeServiceClient {
	conn, err := net.Dial(NET, ADDRESS)
	if err != nil {
		panic(err)
	}
	//	client := thrift.NewClient(thrift.NewFramedReadWriteCloser(conn, 0), thrift.NewBinaryProtocol(true, false))
	client := thrift.NewClient(conn, thrift.NewBinaryProtocol(true, false))

	return TimeServiceClient{client}
}

func LoadTimeSeries() *TimeSeries {
	time := connect()
	timeSeries, err := time.Load(TS_NAME)
	if err != nil {
		panic(err)
	}
	return timeSeries
}
