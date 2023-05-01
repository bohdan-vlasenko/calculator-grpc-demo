package main

import (
	api "calculator-grpc/api/proto"
	"context"
)

type CalculatorGrpcServer struct {
	api.UnimplementedCalculatorServer
}

func (s *CalculatorGrpcServer) mustEmbedUnimplementedCalculatorServer() {
	//TODO implement me
	panic("implement me")
}

func (s *CalculatorGrpcServer) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	return &api.AddResponse{S: req.X + req.Y}, nil
}
