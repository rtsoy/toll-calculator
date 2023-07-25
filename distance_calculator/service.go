package main

import (
	"github.com/rtsoy/toll-calculator/types"
	"math"
)

type CalculatorServicer interface {
	CalculateDistance(data types.OBUData) (float64, error)
}

type CalculatorService struct {
	prevPoint *types.OBUData
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{}
}

func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	var distance float64

	if s.prevPoint != nil {
		distance = calculateDistance(s.prevPoint.Lat, s.prevPoint.Long, data.Lat, data.Long)
	}

	s.prevPoint = &data

	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
