package main

import (
	"fmt"
	"github.com/samuel/go-thrift"
	"net"
	"net/rpc"
	"test"
)

type userStorageServiceImplementation int

func (u *userStorageServiceImplementation) Store(userprofile *test.UserProfile) (bool, error) {
	fmt.Printf("Store user profile with:\n")
	fmt.Printf("\t id: %+v\n", userprofile.Uid)
	fmt.Printf("\t name: %+v\n", userprofile.Name)
	fmt.Printf("\t blurb: %+v\n", userprofile.Blurb)
	return true, nil
}

func main() {
	userStorageService := new(userStorageServiceImplementation)
	rpc.RegisterName("Thrift", &test.UserStorageServer{userStorageService})

	ln, err := net.Listen("tcp", ":1463")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("ERROR: %+v\n", err)
			continue
		}
		fmt.Printf("New connection %+v\n", conn)
		go rpc.ServeCodec(thrift.NewServerCodec(thrift.NewFramedReadWriteCloser(conn, 0), thrift.NewBinaryProtocol(true, false)))
	}
}
