package core

import (
	"golang.org/x/net/context"
	"grpc_micro/contract"
)

func (i *InnerStruct) Check(ctx context.Context, in *contract.Nothing) (*contract.Nothing, error) {
	return new(contract.Nothing), nil
}

func (i *InnerStruct) Add(ctx context.Context, in *contract.Nothing) (*contract.Nothing, error) {
	return new(contract.Nothing), nil
}

func (i *InnerStruct) Test(ctx context.Context, in *contract.Nothing) (*contract.Nothing, error) {
	return new(contract.Nothing), nil
}
