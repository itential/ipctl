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

type Template struct {
	Id            string                 `json:"_id,omitempty"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	Command       string                 `json:"command"`
	Group         string                 `json:"group"`
	Template      string                 `json:"template"`
	Data          string                 `json:"data"`
	Created       string                 `json:"created"`
	CreatedBy     any                    `json:"createdBy"`
	LastUpdated   string                 `json:"lastUpdated"`
	LastUpdatedBy any                    `json:"lastUpdateBy"`
	Namespace     map[string]interface{} `json:"namespace"`
	Tags          []string               `json:"tags"`
}

type TemplateService struct {
	client *ServiceClient
}

func NewTemplateService(iapClient client.Client) *TemplateService {
	return &TemplateService{client: NewServiceClient(iapClient)}
}

func NewTemplate(name, group, description, t string) Template {
	logger.Trace()

	if t == "" {
		t = "textfsm"
	}

	return Template{
		Name:        name,
		Group:       group,
		Description: description,
		Type:        t,
	}
}

// GetAll will return all temlates from the server
func (svc *TemplateService) GetAll() ([]Template, error) {
	logger.Trace()

	var res PaginatedResponse
	var templates []Template

	var limit = 100
	var skip = 0

	// NOTE (privateip) I believe that if the query params are not specified
	// this API will simply return all items which is contrary to the API
	// documentation.  Need to test
	for {
		if err := svc.client.GetRequest(&Request{
			uri:    "/automation-studio/templates",
			params: &QueryParams{Limit: limit, Skip: skip},
		}, &res); err != nil {
			return nil, err
		}

		for _, ele := range res.Items {
			var t Template
			if err := Unmarshal(ele, &t); err != nil {
				return nil, err
			}
			templates = append(templates, t)
		}

		if len(templates) == res.Total {
			break
		}

		skip += limit
	}

	logger.Info("GetAll found %v template(s)", len(templates))

	return templates, nil
}

// Get will return the template specified by its id.
func (svc *TemplateService) Get(id string) (*Template, error) {
	logger.Trace()

	var res *Template
	var uri = fmt.Sprintf("/automation-studio/templates/%s", id)

	// FIXME (privateip) This can be optimzied by using query params instead of
	// iterating over all configured templates
	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetByName will attempt to retrieve the template using the template name.  If
// the template is not found, this function will return an error
func (svc *TemplateService) GetByName(name string) (*Template, error) {
	logger.Trace()

	var res PaginatedResponse

	if err := svc.client.GetRequest(&Request{
		uri: "/automation-studio/templates",
		query: map[string]string{
			"contains[name]": name,
		},
	}, &res); err != nil {
		return nil, err
	}

	if len(res.Items) == 0 {
		return nil, errors.New("template not found")
	}

	var template *Template

	for _, ele := range res.Items {
		if ele.(map[string]interface{})["name"].(string) == name {
			if err := Unmarshal(ele, &template); err != nil {
				return nil, err
			}
			break
		}
	}

	return template, nil
}

func (svc *TemplateService) Create(in Template) (*Template, error) {
	logger.Trace()

	body := map[string]map[string]interface{}{
		"template": map[string]interface{}{
			"name":        in.Name,
			"group":       in.Group,
			"type":        in.Type,
			"description": in.Description,
		},
	}

	type Response struct {
		Template *Template `json:"created"`
		Edit     string    `json:"edit"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/automation-studio/templates",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res.Template, nil
}

func (svc *TemplateService) Delete(id string) error {
	logger.Trace()
	return svc.client.Delete(
		fmt.Sprintf("/automation-studio/templates/%s", id),
	)
}

func (svc *TemplateService) Import(in Template) (*Template, error) {
	logger.Trace()

	body := map[string][]Template{"templates": []Template{in}}

	type Response struct {
		Imported []struct {
			Succcess bool      `json:"success"`
			Message  string    `json:"message"`
			Original *Template `json:"original"`
		} `json:"imported"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/automation-studio/templates/import",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Imported[0].Message)

	return res.Imported[0].Original, nil
}

func (svc *TemplateService) Export(id string) (*Template, error) {
	logger.Trace()

	var res *Template
	var uri = fmt.Sprintf("/automation-studio/templates/%s/export", id)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}
