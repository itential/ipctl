// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type InstanceOperation struct {
	Message  string     `json:"message"`
	Data     []Instance `json:"data"`
	Metadata Metadata   `json:"metadata"`
}

type LastAction struct {
	Id          string `json:"_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	ExecutionId string `json:"executionId"`
}

type Instance struct {
	Id           string                 `json:"_id"`
	ModelId      string                 `json:"modelId"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	InstanceData map[string]interface{} `json:"instanceData"`
	LastAction   LastAction             `json:"lastAction"`
	/*
		Created       string                 `json:"created"`
		CreatedBy     string                 `json:"createdBy"`
		LastUpdated   string                 `json:"lastUpdated"`
		LastUpdatedBy string                 `json:"lastUpdatedBy"`
	*/
}

type InstanceService struct {
	client *client.IapClient
}

func NewInstanceService(iapClient *client.IapClient) *InstanceService {
	return &InstanceService{client: iapClient}
}

func (svc *InstanceService) Get(modelId, instanceId string) (*Instance, error) {
	logger.Trace()

	var response map[string]interface{}
	resp, err := Do(&Request{
		client:   svc.client.Http(),
		method:   http.MethodGet,
		uri:      fmt.Sprintf("/lifecycle-manager/resources/%s/instances/%s", modelId, instanceId),
		response: &response,
	})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(resp.Body))
	}

	var instance *Instance
	if err := Unmarshal(response["data"].(map[string]interface{}), &instance); err != nil {
		return nil, err
	}

	return instance, nil

}

func (svc *InstanceService) GetAll(modelId string) ([]Instance, error) {
	var response InstanceOperation

	Do(&Request{
		client:   svc.client.Http(),
		method:   http.MethodGet,
		uri:      fmt.Sprintf("/lifecycle-manager/resources/%s/instances", modelId),
		response: &response,
	})

	return response.Data, nil
}
