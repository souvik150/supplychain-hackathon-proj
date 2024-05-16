package main

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"

	"github.com/souvik150/supplychain-hackathon-proj/aggregator/client"
	"github.com/souvik150/supplychain-hackathon-proj/types"
)

type DataConsumer interface{
	Start()
}

type KafkaConsumer struct {
	consumer *kafka.Consumer
	isRunning bool
	calcService CalculatorServicer
	aggClient *client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient *client.Client) (DataConsumer, error){
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer: c,
		calcService: svc,
		aggClient: aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("Starting Kafka Transport")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("Kafka consumer error: %v (%v)\n", err, msg)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON Serialization error: %v", err)
			continue
		}

		dist, err := c.calcService.CalculateDistance(data);
		if err != nil {
			logrus.Errorf("Error calculating distance: %v", err)
			continue
		}
		logrus.Infof("Distance calculated: %v", dist)


		req:= types.Distance{
			Value: dist,
			Unix: time.Now().UnixNano(),
			OBUID: data.ODUID,
		}
		if err := c.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("Error aggregating invoice: %v", err)
			continue
		}
	}
}