package aggservice

import (
	"context"
	"github.com/go-kit/log"
	"github.com/rtsoy/toll-calculator/types"
	"time"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	log  log.Logger
	next Service
}

func NewLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			log:  logger,
			next: next,
		}
	}
}

func (lm *loggingMiddleware) Aggregate(ctx context.Context, distance types.Distance) (err error) {
	defer func(start time.Time) {
		lm.log.Log(
			"took", time.Since(start),
			"obuID", distance.OBUID,
			"distance", distance.Value,
			"err", err,
		)
	}(time.Now())

	return lm.next.Aggregate(ctx, distance)
}

func (lm *loggingMiddleware) Calculate(ctx context.Context, id int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		lm.log.Log(
			"took", time.Since(start),
			"id", id,
			"obuID", invoice.OBUID,
			"totalDistance", invoice.TotalDistance,
			"totalAmount", invoice.TotalAmount,
			"err", err,
		)
	}(time.Now())

	return lm.next.Calculate(ctx, id)
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
	return im.next.Aggregate(ctx, distance)
}

func (im instrumentationMiddleware) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	return im.next.Calculate(ctx, id)
}
