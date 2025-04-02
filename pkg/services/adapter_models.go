// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type AdapterModelService struct {
	client *ServiceClient
}

func NewAdapterModelService(c client.Client) *AdapterModelService {
	return &AdapterModelService{client: NewServiceClient(c)}
}

// GetAll will retrieve all of the adapter models that are avalalbe on the
// Itential Platform server and return them as a string array.
func (svc *AdapterModelService) GetAll() ([]string, error) {
	logger.Trace()

	type Response struct {
		Models []string `json:"adapterModels"`
		Total  int      `json:"total"`
	}

	var res Response

	if err := svc.client.Get("/adapter-models/types", &res); err != nil {
		return nil, err
	}

	logger.Info("Found %v adapter model(s)", res.Total)

	return res.Models, nil
}
