package main

import "github.com/sirupsen/logrus"

const kafkaTopic = "obudata"

func main() {
	csv := NewCalculatorService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, csv)
	if err != nil {
		logrus.Fatal(err)
	}

	kafkaConsumer.Start()
}
