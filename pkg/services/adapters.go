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

// AdapterOperationResponse represents the standard response structure
// for adapter operations that return status and message information.
type AdapterOperationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// AdapterProperties contains the core configuration properties for an adapter,
// including its identifier, type, associated brokers, groups, and custom properties.
type AdapterProperties struct {
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Brokers    []string               `json:"brokers"`
	Groups     []any                  `json:"groups"`
	Properties map[string]interface{} `json:"properties"`
}

// Adapter represents an adapter instance configuration within the Itential platform.
// It contains all the necessary information to define, configure, and manage an adapter.
type Adapter struct {
	Name             string                 `json:"name"`
	Type             string                 `json:"type"`
	Model            string                 `json:"model"`
	Properties       AdapterProperties      `json:"properties"`
	IsEncrypted      bool                   `json:"isEncrypted"`
	LoggerProperties map[string]interface{} `json:"loggerProps"`
	Virtual          bool                   `json:"virtual"`
}

// AdapterService provides methods for managing adapter instances in the Itential platform.
// It handles CRUD operations, lifecycle management, and import/export functionality.
type AdapterService struct {
	client *ServiceClient
}

// NewAdapterService creates and returns a new AdapterService instance
// configured with the provided client for API communication.
func NewAdapterService(c client.Client) *AdapterService {
	return &AdapterService{client: NewServiceClient(c)}
}

// GetAll will retrieve all configured adapter instances and return them to the
// calling function as an array of type Adapter.  If there are no configured
// adapters, this function will return an empty array.
func (svc *AdapterService) GetAll() ([]Adapter, error) {
	logger.Trace()

	type Metadata struct {
		ActiveSync bool `json:"activeSync"`
		IsAlive    bool `json:"isAlive"`
	}

	type Results struct {
		Metadata Metadata `json:"metadata"`
		Data     Adapter  `json:"data"`
		Virtual  bool     `json:"virtual"`
	}

	type Collection struct {
		Results []Results `json:"results"`
		Total   int       `json:"total"`
	}

	var res Collection

	if err := svc.client.Get("/adapters", &res); err != nil {
		return nil, err
	}

	var adapters []Adapter

	for _, ele := range res.Results {
		adapters = append(adapters, ele.Data)
	}

	logger.Info("Found %v adapter(s)", res.Total)

	return adapters, nil
}

// Get attempts to retrieve the adapter as specified by the name argument.  If
// the adapter exists, it is returned to the calling function.  If the
// specified adapter does not exist, an error is returned.
func (svc *AdapterService) Get(name string) (*Adapter, error) {
	logger.Trace()

	type Metadata struct {
		ActiveSync bool `json:"activeSync"`
		IsAlive    bool `json:"isAlive"`
	}

	type Response struct {
		Metadata Metadata `json:"metadata"`
		Data     Adapter  `json:"data"`
		Virtual  bool     `json:"virtual"`
	}

	var res Response
	var uri = fmt.Sprintf("/adapters/%s", name)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// Create creates a new adapter instance with the provided configuration.
// It returns the created adapter on success or an error if the operation fails.
func (svc *AdapterService) Create(in Adapter) (*Adapter, error) {
	logger.Trace()

	body := map[string]interface{}{"properties": in}

	type Response struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    *Adapter `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/adapters",
		body:               &body,
		expectedStatusCode: 200,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	return res.Data, nil
}

// Delete removes the adapter instance with the specified name.
// It returns an error if the adapter doesn't exist or the operation fails.
func (svc *AdapterService) Delete(name string) error {
	logger.Trace()
	return svc.client.Delete(fmt.Sprintf("/adapters/%s", name))
}

// Import imports an adapter configuration, creating or updating the adapter instance.
// It returns the imported adapter on success or an error if the operation fails.
func (svc *AdapterService) Import(in Adapter) (*Adapter, error) {
	logger.Trace()

	type Response struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    *Adapter `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/adapters/import",
		body:               map[string]interface{}{"properties": in},
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	return res.Data, nil
}

// Update modifies an existing adapter instance with the provided configuration.
// It returns the updated adapter on success or an error if the operation fails.
func (svc *AdapterService) Update(in Adapter) (*Adapter, error) {
	logger.Trace()

	type Response struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    *Adapter `json:"data"`
	}

	var res Response
	var body = map[string]interface{}{"properties": in}
	var uri = fmt.Sprintf("/adapters/%s", in.Name)

	if err := svc.client.Put(uri, body, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	return res.Data, nil
}

// Export retrieves the adapter configuration for the specified adapter name,
// returning it in a format suitable for backup or import operations.
func (svc *AdapterService) Export(name string) (*Adapter, error) {
	logger.Trace()
	return svc.Get(name)
}

// Start initiates the specified adapter instance, making it active and ready to process requests.
// It returns an error if the adapter doesn't exist or cannot be started.
func (svc *AdapterService) Start(name string) error {
	logger.Trace()
	return svc.client.Put(fmt.Sprintf("/adapters/%s/start", name), nil, nil)
}

// Stop halts the specified adapter instance, making it inactive.
// It returns an error if the adapter doesn't exist or cannot be stopped.
func (svc *AdapterService) Stop(name string) error {
	logger.Trace()
	return svc.client.Put(fmt.Sprintf("/adapters/%s/stop", name), nil, nil)
}

// Restart stops and then starts the specified adapter instance.
// It returns an error if the adapter doesn't exist or if either operation fails.
func (svc *AdapterService) Restart(name string) error {
	logger.Trace()

	if err := svc.Stop(name); err != nil {
		return err
	}

	return svc.Start(name)
	//return svc.client.Put(fmt.Sprintf("/adapters/%s/restart", name), nil, nil)
}
