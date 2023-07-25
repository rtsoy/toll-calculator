package main

import (
	"github.com/rtsoy/toll-calculator/types"
	"github.com/sirupsen/logrus"
	"time"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}

}

func (l *LogMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"distance": dist,
			"err":      err,
			"took":     time.Since(start),
		}).Info("calculate distance")
	}(time.Now())

	return l.next.CalculateDistance(data)
}
