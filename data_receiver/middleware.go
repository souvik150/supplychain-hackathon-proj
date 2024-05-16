package main

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/souvik150/supplychain-hackathon-proj/types"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	
	defer func(start time.Time){
		logrus.WithFields(logrus.Fields{
			"obuID" : data.ODUID,
			"lat": data.Lat,
			"long": data.Long,
			"took": time.Since(start),
			}).Info("Producing to kafka")
	}(time.Now())
			
	return l.next.ProduceData(data)
}