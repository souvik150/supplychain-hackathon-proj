package main

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/souvik150/supplychain-hackathon-proj/types"
)

type DataProducer interface {
	ProduceData (types.OBUData) error
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic string
}

func NewKafkaProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// // Delivery report handler for produced messages
	// go func() {
	// 	for e := range p.Events() {
	// 		switch ev := e.(type) {
	// 		case *kafka.Message:
	// 			if ev.TopicPartition.Error != nil {
	// 				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
	// 			} else {
	// 				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
	// 			}
	// 		}
	// 	}
	// }()

	return &KafkaProducer{
		producer: p,
		topic: topic,
	}, nil
}

func (p *KafkaProducer) ProduceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	partition := int32(0) 
	
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: partition},
		Value: b,
	}, nil)

	return err
}