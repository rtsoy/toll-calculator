package aggservice

import (
	"context"
	"github.com/rtsoy/toll-calculator/types"
)

const basePrice = 44

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type BasicService struct {
	store Storer
}

func newBasicService(store Storer) Service {
	return &BasicService{
		store: store,
	}
}

func (bs *BasicService) Aggregate(ctx context.Context, distance types.Distance) error {
	return bs.store.Insert(distance)
}

func (bs *BasicService) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	dist, err := bs.store.Get(id)
	if err != nil {
		return nil, err
	}

	inv := &types.Invoice{
		OBUID:         id,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}

	return inv, nil
}

func NewAggregatorService(store Storer) Service {
	var svc Service
	{
		svc = newBasicService(store)
		svc = NewLoggingMiddleware()(svc)
		svc = NewInstrumentationMiddleware()(svc)
	}
	return svc
}
