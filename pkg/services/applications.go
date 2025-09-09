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

// ApplicationOperationResponse represents the response structure for application operations
type ApplicationOperationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Application represents an application configuration in the Itential platform
type Application struct {
	Name             string                 `json:"name"`        // Name of the application
	Type             string                 `json:"type"`        // Type of the application
	Model            string                 `json:"model"`       // Model used by the application
	Properties       map[string]interface{} `json:"properties"`  // Application-specific properties
	IsEncrypted      bool                   `json:"isEncrypted"` // Whether the application uses encryption
	LoggerProperties map[string]interface{} `json:"loggerProps"` // Logger configuration properties
}

// ApplicationService provides methods to manage applications in the Itential platform
type ApplicationService struct {
	client *ServiceClient
}

// NewApplicationService creates a new ApplicationService instance with the provided client
func NewApplicationService(c client.Client) *ApplicationService {
	return &ApplicationService{client: NewServiceClient(c)}
}

// GetAll retrieves all applications from the Itential platform
func (svc *ApplicationService) GetAll() ([]Application, error) {
	logger.Trace()

	type Result struct {
		Data     Application `json:"data"`
		Metadata struct {
			ActiveSync bool `json:"activeSync"`
			IsActive   bool `json:"IsActive"`
		} `json:"metadata"`
	}

	type Response struct {
		Results []Result `json:"results"`
		Total   int      `json:"isAlive"`
	}

	var res Response

	if err := svc.client.Get("/applications", &res); err != nil {
		return nil, err
	}

	var values []Application

	for _, ele := range res.Results {
		values = append(values, ele.Data)
	}

	return values, nil
}

// Get retrieves a specific application by name from the Itential platform
func (svc *ApplicationService) Get(name string) (*Application, error) {
	logger.Trace()

	type Response struct {
		Metadata struct {
			ActiveSync bool `json:"activeSync"`
			IsActive   bool `json:"IsActive"`
		} `json:"metadata"`
		Data *Application `json:"data"`
	}

	var res Response
	var uri = fmt.Sprintf("/applications/%s", name)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// Create creates a new application in the Itential platform
func (svc *ApplicationService) Create(in Application) (*Application, error) {
	logger.Trace()

	type Response struct {
		Status  string       `json:"status"`
		Message string       `json:"message"`
		Data    *Application `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/applications",
		body:               map[string]interface{}{"properties": in},
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil
}

// Start starts the specified application by name
func (svc *ApplicationService) Start(name string) error {
	logger.Trace()

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var res Response
	var uri = fmt.Sprintf("/applications/%s/start", name)

	if err := svc.client.Put(uri, nil, &res); err != nil {
		return err
	}

	logger.Info(res.Message)

	return nil
}

// Stop stops the specified application by name
func (svc *ApplicationService) Stop(name string) error {
	logger.Trace()

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var res Response
	var uri = fmt.Sprintf("/applications/%s/stop", name)

	if err := svc.client.Put(uri, nil, &res); err != nil {
		return err
	}

	logger.Info(res.Message)

	return nil
}

// Restart restarts the specified application by name
func (svc *ApplicationService) Restart(name string) error {
	logger.Trace()

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var res Response
	var uri = fmt.Sprintf("/applications/%s/restart", name)

	if err := svc.client.Put(uri, nil, &res); err != nil {
		return err
	}

	return nil
}
