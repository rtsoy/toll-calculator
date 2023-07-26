package client

import (
	"context"
	"github.com/rtsoy/toll-calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Client types.AggregatorClient
}

func NewGRPCClient(endpoint string) (Client, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := types.NewAggregatorClient(conn)

	return &GRPCClient{
		Client: c,
	}, nil
}

func (c *GRPCClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	panic("implement me")
}

func (c *GRPCClient) Aggregate(ctx context.Context, request *types.AggregateRequest) error {
	_, err := c.Client.Aggregate(ctx, request)
	return err
}
