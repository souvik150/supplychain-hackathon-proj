package main

import (
	"github.com/souvik150/supplychain-hackathon-proj/types"
)

const basePrice = 3.15

type Aggregator interface{
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface{
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(id int) (*types.Invoice, error) {
	dist, err := i.store.Get(id)
	if err != nil {
		return nil, err
	}
	
	inv := &types.Invoice{
		OBUID: id,
		TotalDistance: dist,
		Amount: dist * basePrice,
	}

	return inv, nil
}