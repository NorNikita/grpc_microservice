package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_micro/contract"
	"grpc_micro/core"
	"net"
)

func StartMyMicroservice(ctx context.Context, address string, aclData string) error {
	aclMap, errAcl := core.ParseAcl(aclData)
	if errAcl != nil {
		return errAcl
	}

	innerStruct := &core.InnerStruct{
		Acl:                  aclMap,
		LogChannel:           make(chan *contract.Event, 0),
		LogListeners:         make([]*core.ListenerLog, 0),
		StopLogChannel:       make(chan interface{}, 0),
		StatisticChannel:     make(chan *core.Info, 0),
		StatisticListeners:   make([]*core.StatisticListener, 0),
		StopStatisticChannel: make(chan interface{}, 0),
	}

	listen, errListen := net.Listen("tcp", address)
	if errListen != nil {
		return errListen
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(innerStruct.Mono),
		grpc.StreamInterceptor(innerStruct.Flux))

	contract.RegisterAdminServer(server, innerStruct)
	contract.RegisterBizServer(server, innerStruct)

	go innerStruct.SendLogEvent()
	go innerStruct.SendStatisticEvent()

	go func() {
		err := server.Serve(listen)
		if err != nil {
			fmt.Println("ERROR ", err)
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			{
				innerStruct.StopLogChannel <- struct{}{}
				innerStruct.StopStatisticChannel <- struct{}{}
				server.Stop()
				return
			}
		}
	}()

	return nil
}
