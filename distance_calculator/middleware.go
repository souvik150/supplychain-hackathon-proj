package main

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/souvik150/supplychain-hackathon-proj/types"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (l LogMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"dist": dist,
			"err": err,
		}).Info("Distance calculation")
	}(time.Now())
	
	dist, err = l.next.CalculateDistance(data)
	return
}