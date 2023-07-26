package main

import (
	"github.com/rtsoy/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
)

const (
	kafkaTopic  = "obudata"
	aggEndpoint = "http://127.0.0.1:3000"
)

func main() {
	var (
		csv = NewCalculatorService()
	)

	csv = NewLogMiddleware(csv)

	httpClient := client.NewHTTPClient(aggEndpoint)
	// grpcClient, err := client.NewGRPCClient(aggEndpoint)
	//if err != nil {
	//	log.Fatal(err)
	//}

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, csv, httpClient)
	if err != nil {
		logrus.Fatal(err)
	}

	kafkaConsumer.Start()
}
