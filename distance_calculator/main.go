package main

import "github.com/sirupsen/logrus"

const kafkaTopic = "obudata"

func main() {
	var (
		csv  = NewCalculatorService()
		csvm = NewLogMiddleware(csv)
	)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, csvm)
	if err != nil {
		logrus.Fatal(err)
	}

	kafkaConsumer.Start()
}
