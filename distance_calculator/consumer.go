package main

import (
	"encoding/json"
	"time"

	"github.com/0xDarkXnight/Toll-Calculator-Microservice/aggregator/client"
	"github.com/0xDarkXnight/Toll-Calculator-Microservice/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CaclulatorServicer
	aggClient   *client.Client
}

func NewKafkaConsumer(topic string, svc CaclulatorServicer, aggClient *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s", err)
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error: %s", err)
		}
		req := types.Distance{
			Value: distance,
			OBUID: data.OBUID,
			Unix:  time.Now().UnixNano(),
		}
		if err := c.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("aggregate error: %s", err)
			continue
		}
	}
}
