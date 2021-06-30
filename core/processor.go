package core

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"grpc_micro/contract"
	"strings"
	"time"
)

func (i *InnerStruct) Mono(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	metaInfo, _ := metadata.FromIncomingContext(ctx)

	err := i.execute(metaInfo, info.FullMethod)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func (i *InnerStruct) Flux(
	req interface{},
	serverStream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	metaInfo, _ := metadata.FromIncomingContext(serverStream.Context())

	err := i.execute(metaInfo, info.FullMethod)
	if err != nil {
		return err
	}
	return handler(req, serverStream)
}

func (i *InnerStruct) execute(metaInfo metadata.MD, fullMethod string) error {
	consumer, ok := metaInfo["consumer"]
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := i.checkPermissionForMethod(consumer[0], fullMethod)
	if err != nil {
		return err
	}

	inf := &Info{
		method:   fullMethod,
		consumer: consumer[0],
	}
	i.RWMutex.Lock()
	i.StatisticChannel <- inf
	i.RWMutex.Unlock()

	logEvent := &contract.Event{
		Timestamp: time.Now().Unix(),
		Consumer:  consumer[0],
		Method:    fullMethod,
		Host:      strings.Split(metaInfo[":authority"][0], ":")[0] + ":",
	}
	i.RWMutex.Lock()
	i.LogChannel <- logEvent
	i.RWMutex.Unlock()

	return nil
}
