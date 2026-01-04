// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"fmt"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

// IntegrationModelTls represents TLS configuration settings for integration models.
// It defines whether TLS is enabled and certificate validation behavior.
type IntegrationModelTls struct {
	Enabled             bool `json:"enabled"`
	RefjectUnauthorized bool `json:"rejectUnauthroized"`
}

// IntegrationModelServer represents server configuration for integration models.
// It defines the connection details including protocol, host, and base path.
type IntegrationModelServer struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	BasePath string `json:"base_path"`
}

// IntegrationModelProperties represents the configuration properties for an integration model.
// It includes authentication, server settings, TLS configuration, and version information.
type IntegrationModelProperties struct {
	Authentication map[string]interface{} `json:"authentication"`
	Server         IntegrationModelServer `json:"server"`
	Tls            IntegrationModelTls    `json:"tls"`
	Version        string                 `json:"version"`
}

// IntegrationModel represents a model definition for integrations in the Itential Platform.
// Models define the structure and configuration templates for creating integrations.
type IntegrationModel struct {
	Model       string                     `json:"model,omitempty"`
	VersionId   string                     `json:"versionId"`
	Description string                     `json:"description"`
	Properties  IntegrationModelProperties `json:"properties"`
}

// IntegrationModelService provides methods for managing integration models in the Itential Platform.
// It handles CRUD operations and export functionality for integration model definitions.
type IntegrationModelService struct {
	BaseService
}

// NewIntegrationModelService creates a new IntegrationModelService instance with the provided HTTP client.
// The client is used to communicate with the Itential Platform integration models API.
func NewIntegrationModelService(c client.Client) *IntegrationModelService {
	return &IntegrationModelService{BaseService: NewBaseService(c)}
}

// GetAll retrieves all integration models from the Itential Platform.
// It sends a GET request to /integration-models and returns the complete list of available models.
// Returns a slice of all integration models or an error if the operation fails.
func (svc *IntegrationModelService) GetAll() ([]IntegrationModel, error) {
	logger.Trace()

	type Response struct {
		IntegrationModels []IntegrationModel `json:"integrationModels"`
		Total             int                `json:"total"`
	}

	var res Response

	if err := svc.BaseService.Get("/integration-models", &res); err != nil {
		return nil, err
	}

	return res.IntegrationModels, nil
}

// Get retrieves a specific integration model by its name from the Itential Platform.
// It sends a GET request to /integration-models/{name}.
// Returns the integration model definition or an error if the operation fails or model is not found.
func (svc *IntegrationModelService) Get(name string) (*IntegrationModel, error) {
	logger.Trace()

	var res *IntegrationModel
	var uri = fmt.Sprintf("/integration-models/%s", name)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// Create creates a new integration model in the Itential Platform.
// It sends a POST request to /integration-models with the model definition.
// After creation, it retrieves and returns the created model using the returned versionId.
// Returns the created integration model or an error if the operation fails.
func (svc *IntegrationModelService) Create(in map[string]interface{}) (*IntegrationModel, error) {
	logger.Trace()

	body := map[string]interface{}{"model": in}

	type Response struct {
		Message string                 `json:"message"`
		Status  string                 `json:"status"`
		Data    map[string]interface{} `json:"data"`
	}

	var res Response

	if err := svc.PostRequest(&Request{
		uri:                "/integration-models",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	model, err := svc.Get(res.Data["versionId"].(string))
	if err != nil {
		return nil, err
	}

	return model, nil
}

// Delete removes an integration model from the Itential Platform by its name.
// It sends a DELETE request to /integration-models/{name}.
// Returns an error if the operation fails or the model is not found.
func (svc *IntegrationModelService) Delete(name string) error {
	logger.Trace()
	return svc.BaseService.Delete(fmt.Sprintf("/integration-models/%s", name))
}

// Export retrieves the exportable definition of an integration model by its name.
// It sends a GET request to /integration-models/{name}/export.
// Returns a map containing the model's exportable configuration or an error if the operation fails.
func (svc *IntegrationModelService) Export(name string) (map[string]interface{}, error) {
	logger.Trace()

	var res map[string]interface{}
	var uri = fmt.Sprintf("/integration-models/%s/export", name)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}
