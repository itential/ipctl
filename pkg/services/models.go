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

type RunActionRequest struct {
	ActionId     string                 `json:"actionId"`
	Inputs       map[string]interface{} `json:"inputs"`
	Name         string                 `json:"name"`
	Instance     string                 `json:"instance"`
	InstanceName string                 `json:"instanceName"`
}

type RunActionResponse struct {
	Id                  string                 `json:"_id"`
	ModelId             string                 `json:"modelId"`
	InstanceId          string                 `json:"instanceId"`
	ActionId            string                 `json:"actionId"`
	ModelName           string                 `json:"modelName"`
	InstanceName        string                 `json:"instanceName"`
	ActionName          string                 `json:"actionName"`
	StartTime           string                 `json:"startTime"`
	EndTime             string                 `json:"endTime"`
	Initiator           string                 `json:"initiator"`
	JobId               string                 `json:"jobId"`
	Status              string                 `json:"status"`
	Progress            map[string]interface{} `json:"progress"`
	Errors              []string               `json:"errors"`
	InitialInstanceData map[string]interface{} `json:"initialInstanceData"`
	FinalInstanceData   map[string]interface{} `json:"finalInstanceData"`
}

type ModelOperation struct {
	Message  string   `json:"message"`
	Data     []Model  `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type ModelAction struct {
	Id              string `json:"_id"`
	Name            string `json:"name"`
	PreWorkflowJst  string `json:"preWorkflowJst,omitempty"`
	PostWorkflowJst string `json:"postWorkflowJst,omitempty"`
	Workflow        string `json:"workflow,omitempty"`
	Type            string `json:"type,omitempty"`
}

type Model struct {
	Id            string                 `json:"_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Schema        map[string]interface{} `json:"schema"`
	Actions       []ModelAction          `json:"actions"`
	Created       string                 `json:"created"`
	CreatedBy     any                    `json:"createdBy"`
	LastUpdated   string                 `json:"lastUpdated"`
	LastUpdatedBy any                    `json:"lastUpdatedBy"`
}

type ModelService struct {
	client *ServiceClient
}

func NewModel(name, desc string) Model {
	logger.Trace()

	return Model{
		Name:        name,
		Description: desc,
	}
}

func NewModelService(c client.Client) *ModelService {
	return &ModelService{client: NewServiceClient(c)}
}

func (svc *ModelService) Get(id string) (*Model, error) {
	logger.Trace()

	var res *Model
	var uri = fmt.Sprintf("/lifecycle-manager/resources/%s", id)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *ModelService) GetAll() ([]Model, error) {
	logger.Trace()

	var res ModelOperation
	var uri = "/lifecycle-manager/resources"

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (svc *ModelService) GetByName(name string) (*Model, error) {
	logger.Trace()

	models, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	var m *Model

	for _, ele := range models {
		if ele.Name == name {
			m = &ele
			break
		}
	}

	if m == nil {
		return nil, errors.New("model not found")
	}

	return m, nil
}

func (svc *ModelService) Create(in Model) (*Model, error) {
	logger.Trace()

	body := map[string]interface{}{
		"name":        in.Name,
		"description": in.Description,
	}

	if in.Schema != nil {
		in.Schema["$id"] = in.Name
		body["schema"] = in.Schema
	}

	type Response struct {
		Message  string                 `json:"message"`
		Data     *Model                 `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/lifecycle-manager/resources",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (svc *ModelService) Delete(id string, deleteInstances bool) error {
	logger.Trace()

	req := &Request{
		uri: fmt.Sprintf("/lifecycle-manager/resources/%s", id),
	}

	if deleteInstances {
		req.body = map[string]interface{}{
			"queryParameters": map[string]interface{}{
				"delete-associated-instances": "true",
			},
		}
	}

	return svc.client.DeleteRequest(req, nil)
}

func (svc *ModelService) RunAction(model string, in RunActionRequest) (*RunActionResponse, error) {
	logger.Trace()

	var res *RunActionResponse
	var uri = fmt.Sprintf("/lifecycle-manager/resources/%s/run-action", model)

	if err := svc.client.Post(uri, &in, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *ModelService) Import(in Model) (*Model, error) {
	logger.Trace()

	body := map[string]interface{}{
		"model": in,
	}

	type Response struct {
		Message string `json:"message"`
		Data    *Model `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/lifecycle-manager/resources/import",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (svc *ModelService) Export(id string) (*Model, error) {
	logger.Trace()

	type Response struct {
		Message string `json:"message"`
		Data    *Model `json:"data"`
	}

	var res Response

	var uri = fmt.Sprintf("/lifecycle-manager/resources/%s/export", id)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil
}
