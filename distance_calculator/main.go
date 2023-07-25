package main

import (
	"github.com/rtsoy/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
)

const (
	kafkaTopic  = "obudata"
	aggEndpoint = "http://127.0.0.1:3000/aggregate"
)

func main() {
	var (
		csv    = NewCalculatorService()
		client = client.NewClient(aggEndpoint)
	)

	csv = NewLogMiddleware(csv)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, csv, client)
	if err != nil {
		logrus.Fatal(err)
	}

	kafkaConsumer.Start()
}
