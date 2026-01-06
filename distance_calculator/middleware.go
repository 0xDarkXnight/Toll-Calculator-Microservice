package main

import (
	"time"

	"github.com/0xDarkXnight/Toll-Calculator-Microservice/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CaclulatorServicer
}

func NewLogMiddleware(next CaclulatorServicer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"dist": dist,
			"err":  err,
			"took": time.Since(start),
		}).Info("calculate distance")
	}(time.Now())
	dist, err = m.next.CalculateDistance(data)
	return
}
