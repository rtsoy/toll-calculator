package aggendpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/rtsoy/toll-calculator/go-kit-example/aggservice"
	"github.com/rtsoy/toll-calculator/types"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
	Err           error   `json:"err"`
}

type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type AggregateResponse struct {
	Err error `json:"err"`
}

func (s Set) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, CalculateRequest{
		OBUID: id,
	})
	if err != nil {
		return nil, err
	}

	result := resp.(CalculateResponse)

	return &types.Invoice{
		OBUID:         result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount:   result.TotalAmount,
	}, nil
}

func (s Set) Aggregate(ctx context.Context, distance types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		Value: distance.Value,
		OBUID: distance.OBUID,
		Unix:  distance.Unix,
	})
	return err
}

func MakeAggregateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, types.Distance{
			Value: req.Value,
			OBUID: req.OBUID,
			Unix:  req.Unix,
		})

		return AggregateResponse{Err: err}, err
	}
}

func MakeCalculateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		v, err := s.Calculate(ctx, req.OBUID)

		return CalculateResponse{
			OBUID:         v.OBUID,
			TotalDistance: v.TotalDistance,
			TotalAmount:   v.TotalAmount,
			Err:           err,
		}, err
	}
}
