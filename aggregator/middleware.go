package main

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/souvik150/supplychain-hackathon-proj/types"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (l LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"distance": distance,
			"err": err,
		}).Info("Aggregating distance")
	}(time.Now())
	
	err = l.next.AggregateDistance(distance)
	return err
}

func (l LogMiddleware) CalculateInvoice(id int) (inv *types.Invoice, err error) {
	var (
		distance float64
		amount float64
	)
	if inv != nil {
		distance = inv.TotalDistance
		amount = inv.Amount
	}

	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"OBUID": id,
			"err": err,
			"distance": distance,
			"amount": amount,
		}).Info("Calculating invoice")
	}(time.Now())
	
	inv, err = l.next.CalculateInvoice(id)
	return inv, err
}