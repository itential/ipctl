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

// RunActionRequest represents the payload for executing a model action
type RunActionRequest struct {
	ActionId     string                 `json:"actionId"`
	Inputs       map[string]interface{} `json:"inputs"`
	Name         string                 `json:"name"`
	Instance     string                 `json:"instance"`
	InstanceName string                 `json:"instanceName"`
}

// RunActionResponse represents the response from executing a model action
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

// ModelOperation represents the response structure for model operations that return multiple models
type ModelOperation struct {
	Message  string   `json:"message"`
	Data     []Model  `json:"data"`
	Metadata Metadata `json:"metadata"`
}

// ModelAction represents an action that can be performed on a model
type ModelAction struct {
	Id              string `json:"_id"`
	Name            string `json:"name"`
	PreWorkflowJst  string `json:"preWorkflowJst,omitempty"`
	PostWorkflowJst string `json:"postWorkflowJst,omitempty"`
	Workflow        string `json:"workflow,omitempty"`
	Type            string `json:"type,omitempty"`
}

// Model represents a lifecycle manager resource model
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

// ModelService provides operations for managing lifecycle manager models
type ModelService struct {
	BaseService
}

// NewModel creates a new Model instance with the specified name and description
func NewModel(name, desc string) Model {
	logger.Trace()

	return Model{
		Name:        name,
		Description: desc,
	}
}

// NewModelService creates a new ModelService instance with the provided client
func NewModelService(c client.Client) *ModelService {
	return &ModelService{BaseService: NewBaseService(c)}
}

// Get retrieves a model by its ID from the lifecycle manager
func (svc *ModelService) Get(id string) (*Model, error) {
	logger.Trace()

	var res *Model
	var uri = fmt.Sprintf("/lifecycle-manager/resources/%s", id)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetAll retrieves all models from the lifecycle manager
func (svc *ModelService) GetAll() ([]Model, error) {
	logger.Trace()

	var res ModelOperation
	var uri = "/lifecycle-manager/resources"

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// GetByName retrieves a model by name using client-side filtering.
// DEPRECATED: Business logic method - prefer using resources.ModelResource.GetByName
func (svc *ModelService) GetByName(name string) (*Model, error) {
	logger.Trace()

	models, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	for i := range models {
		if models[i].Name == name {
			return &models[i], nil
		}
	}

	return nil, errors.New("model not found")
}

// Create creates a new model in the lifecycle manager
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

	if err := svc.PostRequest(&Request{
		uri:                "/lifecycle-manager/resources",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// Delete removes a model from the lifecycle manager
// If deleteInstances is true, associated instances will also be deleted
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

	return svc.DeleteRequest(req, nil)
}

// RunAction executes an action on a model instance
func (svc *ModelService) RunAction(model string, in RunActionRequest) (*RunActionResponse, error) {
	logger.Trace()

	var res *RunActionResponse
	var uri = fmt.Sprintf("/lifecycle-manager/resources/%s/run-action", model)

	if err := svc.Post(uri, &in, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// Import imports a model into the lifecycle manager
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

	if err := svc.PostRequest(&Request{
		uri:                "/lifecycle-manager/resources/import",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// Export exports a model from the lifecycle manager by its ID
func (svc *ModelService) Export(id string) (*Model, error) {
	logger.Trace()

	type Response struct {
		Message string `json:"message"`
		Data    *Model `json:"data"`
	}

	var res Response

	var uri = fmt.Sprintf("/lifecycle-manager/resources/%s/export", id)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	return res.Data, nil
}
