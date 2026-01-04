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

// JsonForm represents a JSON form asset in the Itential Platform.
// JSON forms provide a way to dynamically generate user interfaces
// for data collection based on JSON schema definitions.
type JsonForm struct {
	Id               string                 `json:"id"`               // Unique identifier for the JSON form
	Name             string                 `json:"name"`             // Human-readable name of the form
	Description      string                 `json:"description"`      // Optional description of the form's purpose
	Struct           map[string]interface{} `json:"struct"`           // Form structure definition
	Schema           map[string]interface{} `json:"schema"`           // JSON schema defining form data structure
	UISchema         map[string]interface{} `json:"uiSchema"`         // UI layout and presentation schema
	BindingSchema    map[string]interface{} `json:"bindingSchema"`    // Data binding configuration
	ValidationSchema map[string]interface{} `json:"validationSchema"` // Validation rules for form data
	Created          string                 `json:"created"`          // Creation timestamp in ISO format
	CreatedBy        string                 `json:"createdBy"`        // User ID who created the form
	LastUpdated      string                 `json:"lastUpdated"`      // Last modification timestamp in ISO format
	LastUpdatedBy    string                 `json:"lastUpdatedBy"`    // User ID who last modified the form
	Version          string                 `json:"version"`          // Platform version compatibility
}

// JsonFormService provides methods for managing JSON Form assets
// in the Itential Platform. It handles CRUD operations and bulk
// management of forms used for dynamic UI generation.
type JsonFormService struct {
	BaseService
}

// NewJsonFormService creates and returns a new JsonFormService instance
// configured with the provided client for API communication.
func NewJsonFormService(c client.Client) *JsonFormService {
	return &JsonFormService{BaseService: NewBaseService(c)}
}

// NewJsonForm creates and returns a new JsonForm instance with all fields
// initialized to sensible defaults. The returned object can be passed
// unchanged to JsonFormService.Create() or customized before creation.
// The name and description parameters are required and will be used in
// both the form metadata and the generated JSON schema.
func NewJsonForm(name, description string) JsonForm {
	logger.Trace()

	jf := JsonForm{Name: name, Description: description}

	jf.Schema = map[string]interface{}{
		"description": description,
		"properties":  map[string]interface{}{},
		"required":    []any{},
		"title":       name,
		"type":        "object",
	}

	jf.Struct = map[string]interface{}{
		"description": "",
		"items":       []any{},
		"type":        "object",
	}

	jf.UISchema = map[string]interface{}{}
	jf.BindingSchema = map[string]interface{}{}
	jf.ValidationSchema = map[string]interface{}{}

	return jf

}

// GetAll retrieves all JSON Form assets configured on the server.
// Returns an empty slice if no forms are configured. This method
// does not perform any filtering or sorting of the results.
func (svc *JsonFormService) GetAll() ([]JsonForm, error) {
	logger.Trace()

	var res []JsonForm

	if err := svc.BaseService.Get("/json-forms/forms", &res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get attempts to retrieve the JSON Form asset identified by the id argument.
// If the asset exists on the server, it is returned to the calling function.
// If the asset does not exist, an error is returned.
func (svc *JsonFormService) Get(id string) (*JsonForm, error) {
	logger.Trace()

	var res *JsonForm
	var uri = fmt.Sprintf("/json-forms/forms/%s", id)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("jsonform not found")
	}

	return res, nil
}

// Create creates a new JSON Form asset in the Itential Platform.
// The input JsonForm must have at minimum a valid Name field.
// For forms with proper defaults, use NewJsonForm() to create the input.
// Returns the created form with server-generated fields populated,
// or an error if the creation fails.
func (svc *JsonFormService) Create(in JsonForm) (*JsonForm, error) {
	logger.Trace()

	type Response struct {
		Status  string    `json:"status"`
		Doc     *JsonForm `json:"doc"`
		Message string    `json:"message"`
	}

	var res Response

	if err := svc.PostRequest(&Request{
		uri:                "/json-forms/forms",
		body:               &in,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	if res.Status == "failure" {
		return nil, errors.New(res.Message)
	}

	return res.Doc, nil
}

// Delete removes one or more JSON Form assets from the server by their IDs.
// This is a bulk operation that can delete multiple forms in a single request.
// Returns an error if any of the deletions fail.
func (svc *JsonFormService) Delete(ids []string) error {
	logger.Trace()
	return svc.DeleteRequest(&Request{
		uri:  "/json-forms/forms",
		body: map[string]interface{}{"ids": ids},
	}, nil)
}

// GetByName retrieves a JSON Form asset by its name field.
// This method performs a client-side search through all forms,
// so it may be less efficient than Get() for large numbers of forms.
// Returns an error if no form with the specified name is found.
func (svc *JsonFormService) GetByName(name string) (*JsonForm, error) {
	logger.Trace()

	elements, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	var res *JsonForm

	for _, ele := range elements {
		if ele.Name == name {
			res = &ele
			break
		}
	}

	if res == nil {
		return nil, errors.New("jsonform not found")
	}

	return res, nil
}

// Clear removes all JSON Form assets from the server.
// This is a destructive operation that cannot be undone.
// Use with caution in production environments.
func (svc *JsonFormService) Clear() error {
	logger.Trace()

	elements, err := svc.GetAll()
	if err != nil {
		return err
	}

	var ids []string

	for _, ele := range elements {
		ids = append(ids, ele.Id)
	}

	if len(ids) > 0 {
		if err := svc.Delete(ids); err != nil {
			return err
		}
	}

	return nil
}

// Import imports a JSON Form asset into the Itential Platform using
// the platform's import API. This method can handle forms exported
// from other platform instances. Returns the imported form with
// server-generated fields populated.
func (svc *JsonFormService) Import(in JsonForm) (*JsonForm, error) {
	logger.Trace()

	body := map[string]interface{}{
		"forms": []JsonForm{in},
	}

	type ImportedResponse struct {
		Success  bool                   `json:"success"`
		Message  string                 `json:"message"`
		Original map[string]interface{} `json:"original"`
		Created  map[string]interface{} `json:"created"`
	}

	type Response struct {
		Message  string             `json:"message"`
		Status   string             `json:"status"`
		Imported []ImportedResponse `json:"imported"`
	}

	var res Response

	if err := svc.PostRequest(&Request{
		uri:                "/json-forms/import/forms",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Imported[0].Message)

	jf, err := svc.Get(res.Imported[0].Created["_id"].(string))
	if err != nil {
		return nil, err
	}

	return jf, nil
}
