package main

import (
	"github.com/rtsoy/toll-calculator/types"
	"github.com/sirupsen/logrus"
	"time"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) CalculateInvoice(id int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"took":  time.Since(start),
		}).Info("calculate invoice")
	}(time.Now())

	return l.next.CalculateInvoice(id)
}

func (l *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"id":    distance.OBUID,
			"error": err,
			"took":  time.Since(start),
		}).Info("aggregate distance")
	}(time.Now())

	return l.next.AggregateDistance(distance)
}
