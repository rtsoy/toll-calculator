package aggservice

import (
	"context"
	"github.com/rtsoy/toll-calculator/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next Service
}

func NewLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next: next,
		}
	}
}

func (lm *loggingMiddleware) Aggregate(ctx context.Context, distance types.Distance) error {
	return nil
}

func (lm *loggingMiddleware) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	return nil, nil
}

type instrumentationMiddleware struct {
	next Service
}

func NewInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return &instrumentationMiddleware{
			next: next,
		}
	}
}

func (im instrumentationMiddleware) Aggregate(ctx context.Context, distance types.Distance) error {
	return nil
}

func (im instrumentationMiddleware) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	return nil, nil
}
