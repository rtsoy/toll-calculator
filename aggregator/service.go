package main

import (
	"github.com/rtsoy/toll-calculator/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) CalculateInvoice(id int) (*types.Invoice, error) {
	dist, err := i.store.Get(id)
	if err != nil {
		return nil, err
	}

	invoice := &types.Invoice{
		OBUID:         id,
		TotalDistance: dist,
		TotalAmount:   44,
	}

	return invoice, nil
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	return i.store.Insert(distance)
}
