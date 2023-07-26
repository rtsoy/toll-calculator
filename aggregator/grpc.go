package main

import (
	"context"
	"github.com/rtsoy/toll-calculator/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, request *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		Value: request.Value,
		OBUID: int(request.ObuID),
		Unix:  request.Unix,
	}

	return &types.None{}, s.svc.AggregateDistance(distance)
}
