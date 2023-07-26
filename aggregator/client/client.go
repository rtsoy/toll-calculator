package client

import (
	"context"
	"github.com/rtsoy/toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
