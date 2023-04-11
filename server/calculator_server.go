package server

import (
	"context"
	api "google.golang.org/grpc/examples/calculator/api/proto"
)

type CalculatorGrpcServer struct {
	api.UnimplementedCalculatorServer
}

func (s *CalculatorGrpcServer) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	return &api.AddResponse{S: req.X + req.Y}, nil
}
