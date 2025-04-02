// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type PrometheusMetric struct {
	Name  string
	Type  string
	Value int
	Help  string
}

type MetricService struct {
	client *ServiceClient
}

func NewMetricService(c client.Client) *MetricService {
	return &MetricService{client: NewServiceClient(c)}
}

// Get will retrieve the server Prometheus metrics and return them to to
// calling function.  The format for the return is a string in text format.
func (svc *MetricService) Get() string {
	logger.Trace()

	if err := svc.client.Get("/prometheus_metrics", nil); err != nil {
		return ""
	}

	return ""
}
