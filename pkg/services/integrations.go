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

// Integration represents an integration configuration in the Itential Platform.
// Integrations connect external systems and services to the platform for automation workflows.
type Integration struct {
	Name               string                 `json:"name"`
	Type               string                 `json:"type"`
	Model              string                 `json:"model"`
	Properties         map[string]interface{} `json:"properties"`
	IsEncrypted        bool                   `json:"isEncrypted"`
	LoggerProperties   map[string]interface{} `json:"loggerProps"`
	Virtual            bool                   `json:"virtual"`
	Profiling          bool                   `json:"profiling"`
	SystemProperties   map[string]interface{} `json:"systemProps"`
	EventDeduplciation map[string]interface{} `json:"eventDeduplication"`
}

// IntegrationService provides methods for managing integrations in the Itential Platform.
// It handles CRUD operations for integration configurations and settings.
type IntegrationService struct {
	BaseService
}

// NewIntegrationService creates a new IntegrationService instance with the provided HTTP client.
// The client is used to communicate with the Itential Platform integrations API.
func NewIntegrationService(c client.Client) *IntegrationService {
	return &IntegrationService{BaseService: NewBaseService(c)}
}

// NewIntegration creates a new Integration instance with the specified name and type.
// It initializes the integration with default properties including ID and type fields.
// This is a helper function for creating integration configurations programmatically.
func NewIntegration(name, integrationType string) Integration {
	logger.Trace()

	return Integration{
		Name: name,
		Properties: map[string]interface{}{
			"id":   name,
			"type": integrationType,
		},
	}
}

// Create creates a new integration in the Itential Platform.
// It sends a POST request to /integrations with the integration configuration.
// The Type field is automatically set to "Adapter" and Virtual to true as required by the API.
// Returns the created integration or an error if the operation fails.
func (svc *IntegrationService) Create(in Integration) (*Integration, error) {
	logger.Trace()

	// Make sure to set the Type and Virtual fields to these values otherwise
	// the POST call will return an error
	in.Type = "Adapter"
	in.Virtual = true

	type Response struct {
		Status  string       `json:"status"`
		Message string       `json:"message"`
		Data    *Integration `json:"data"`
	}

	var res Response

	if err := svc.PostRequest(&Request{
		uri:                "/integrations",
		body:               map[string]interface{}{"properties": in},
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	return res.Data, nil

}

// Delete removes an integration from the Itential Platform by its name.
// It sends a DELETE request to /integrations/{name}.
// Returns an error if the operation fails or the integration is not found.
func (svc *IntegrationService) Delete(name string) error {
	logger.Trace()
	return svc.BaseService.Delete(fmt.Sprintf("/integrations/%s", name))
}

// Get retrieves a specific integration by its name from the Itential Platform.
// It sends a GET request to /integrations/{name}.
// Returns the integration configuration or an error if the operation fails or integration is not found.
func (svc *IntegrationService) Get(name string) (*Integration, error) {
	logger.Trace()

	type Response struct {
		Data     *Integration           `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	var res Response
	var uri = fmt.Sprintf("/integrations/%s", name)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// GetAll retrieves all integrations from the Itential Platform.
// It sends a GET request to /integrations and processes the paginated results.
// Returns a slice of all integration configurations or an error if the operation fails.
func (svc *IntegrationService) GetAll() ([]Integration, error) {
	logger.Trace()

	type Result struct {
		Data     *Integration           `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	type Response struct {
		Results []Result `json:"results"`
		Total   int      `json:"total"`
	}

	var res Response

	if err := svc.BaseService.Get("/integrations", &res); err != nil {
		return nil, err
	}

	var elements []Integration

	for _, ele := range res.Results {
		elements = append(elements, *ele.Data)
	}

	return elements, nil
}
