package main

import (
	"log"

	"github.com/0xDarkXnight/Toll-Calculator-Microservice/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"
)

// Transport (HTTP, GRPC, Kafka) -> attach business logic to this transport

func main() {
	var (
		err     error
		service CaclulatorServicer
	)
	service = NewCalculatorService()
	service = NewLogMiddleware(service)
	client := client.NewClient(aggregatorEndpoint)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, service, client)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
