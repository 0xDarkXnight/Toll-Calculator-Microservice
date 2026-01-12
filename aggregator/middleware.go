package main

import (
	"time"

	"github.com/0xDarkXnight/Toll-Calculator-Microservice/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"took": time.Since(start),
		}).Info("Aggregate Distance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}
