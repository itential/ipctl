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

type JsonForm struct {
	Id               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Struct           map[string]interface{} `json:"struct"`
	Schema           map[string]interface{} `json:"schema"`
	UISchema         map[string]interface{} `json:"uiSchema"`
	BindingSchema    map[string]interface{} `json:"bindingSchema"`
	ValidationSchema map[string]interface{} `json:"validationSchema"`
	Created          string                 `json:"created"`
	CreatedBy        string                 `json:"createdBy"`
	LastUpdated      string                 `json:"lastUpdated"`
	LastUpdatedBy    string                 `json:"lastUpdatedBy"`
	Version          string                 `json:"version"`
}

type JsonFormService struct {
	client *ServiceClient
}

func NewJsonFormService(iapClient client.Client) *JsonFormService {
	return &JsonFormService{client: NewServiceClient(iapClient)}
}

// NewJsonForm will create and return an instance of JsorForm with all fields
// set to defaults.  The returned object can be passed unchanged into the
// JsonFormService.Create function.
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

// GetAll will retrieve all JSON Forms assests configured on the server and
// return them as an array.  If there are no cofigured JSON Forms on the
// server, this function will return an empty array.
func (svc *JsonFormService) GetAll() ([]JsonForm, error) {
	logger.Trace()

	var res []JsonForm

	if err := svc.client.Get("/json-forms/forms", &res); err != nil {
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

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("jsonform not found")
	}

	return res, nil
}

// Create will create a new JSON form asset in Platform.  It accepts one
// requirement argument `in`.  The argument must be an instance of JsonForm.
// In order to create a JSON form with all defaults properly set, use the
// `NewJsonForm` function.
func (svc *JsonFormService) Create(in JsonForm) (*JsonForm, error) {
	logger.Trace()

	type Response struct {
		Status  string    `json:"status"`
		Doc     *JsonForm `json:"doc"`
		Message string    `json:"message"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
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

func (svc *JsonFormService) Delete(ids []string) error {
	logger.Trace()
	return svc.client.DeleteRequest(&Request{
		uri:  "/json-forms/forms",
		body: map[string]interface{}{"ids": ids},
	}, nil)
}

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

// Import will call the Itential Platform import API and import the specified
// JSON Form to the server.
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
		Message  string           `json:"message"`
		Status   string           `json:"status"`
		Imported ImportedResponse `json:"imported"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/json-forms/import/forms",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Imported.Message)

	jf, err := svc.Get(res.Imported.Created["_id"].(string))
	if err != nil {
		return nil, err
	}

	return jf, nil
}
