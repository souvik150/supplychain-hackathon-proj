package main

import (
	"log"

	"github.com/souvik150/supplychain-hackathon-proj/aggregator/client"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

const(
	kafkaTopic = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"
)

func main(){
	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient(
		aggregatorEndpoint,
	))
	
	if err != nil {
		log.Fatalf("Error creating Kafka Consumer: %v", err)
	}

	kafkaConsumer.Start()
}