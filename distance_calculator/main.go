package main

import (
	"log"
)

const kafkaTopic = "obudata"

// Transport (HTTP, GRPC, Kafka) -> attach business logic to this transport

func main() {
	var (
		err     error
		service CaclulatorServicer
	)
	service = NewCalculatorService()
	service = NewLogMiddleware(service)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, service)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
