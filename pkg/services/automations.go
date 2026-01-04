// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

// automationCollection represents a collection of automations returned by the API
type automationCollection struct {
	Message  string       `json:"message"`
	Data     []Automation `json:"data"`
	Metadata Metadata     `json:"metadata"`
}

// AutomationGbacEntry represents a GBAC (Group-Based Access Control) entry for automations
type AutomationGbacEntry struct {
	Name        string `json:"name"`
	Provenance  string `json:"provenance"`
	Description string `json:"description"`
}

// AutomationGbac represents GBAC permissions for an automation
type AutomationGbac struct {
	Write []interface{} `json:"write"`
	Read  []interface{} `json:"read"`
}

// Automation represents an automation in the Itential Platform
type Automation struct {
	Id            string `json:"_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ComponentName string `json:"componentName"`
	ComponentType string `json:"componentType"`

	// ComponentId field does not exist when exporting the automation but it
	// does exist when getting it
	ComponentId string `json:"componentId,omitempty"`

	Gbac          AutomationGbac `json:"gbac"`
	Created       string         `json:"created"`
	CreatedBy     string         `json:"createdBy"`
	LastUpdated   string         `json:"lastUpdated"`
	LastUpdatedBy string         `json:"lastUpdatedBy"`

	// Triggers does not exist when getting the autoatmion but it does exists
	// when exporting it
	Triggers []Trigger `json:"triggers,omitempty"`
}

// AutomationService provides methods for managing automations
type AutomationService struct {
	BaseService
}

// NewAutomation creates a new Automation instance with the given name and description
func NewAutomation(name, desc string) Automation {
	logger.Trace()
	return Automation{
		Name:          name,
		Description:   desc,
		ComponentType: "workflows",
	}
}

// NewAutomationService creates a new AutomationService with the given client
func NewAutomationService(c client.Client) *AutomationService {
	return &AutomationService{BaseService: NewBaseService(c)}
}

// Get implements `GET /operations-manager/automations/{id}`
func (svc *AutomationService) Get(id string) (*Automation, error) {
	logger.Trace()

	type Response struct {
		Message string      `json:"message"`
		Data    *Automation `json:"data"`
	}

	var res Response
	var uri = fmt.Sprintf("/operations-manager/automations/%s", id)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	return res.Data, nil
}

// GetByName retrieves an automation by its name
func (svc *AutomationService) GetByName(name string) (*Automation, error) {
	logger.Trace()

	automations, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	var automation *Automation

	for _, ele := range automations {
		if ele.Name == name {
			automation = ele
			break
		}
	}

	if automation == nil {
		return nil, errors.New("automation not found")
	}

	return automation, nil

}

// Create implements `POST /operations-manager/automations`
func (svc *AutomationService) Create(in Automation) (*Automation, error) {
	logger.Trace()

	body := map[string]interface{}{
		"name":          in.Name,
		"description":   in.Description,
		"componentType": in.ComponentType,
	}

	type Response struct {
		Message  string                 `json:"message"`
		Data     *Automation            `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	var res Response

	if err := svc.PostRequest(&Request{
		uri:                "/operations-manager/automations",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	return res.Data, nil
}

// Delete implements `DELETE /operations-manager/automations/{id}`
func (svc *AutomationService) Delete(id string) error {
	logger.Trace()
	return svc.Delete(fmt.Sprintf("/operations-manager/automations/%s", id))
}

// GetAll implements `GET /operations-manager/automations`
func (svc *AutomationService) GetAll() ([]*Automation, error) {
	logger.Trace()

	type Response struct {
		Message  string        `json:"message"`
		Data     []*Automation `json:"data"`
		Metadata Metadata      `json:"metadata"`
	}

	var res Response

	var automations []*Automation

	var limit = 100
	var skip = 0

	for {
		if err := svc.GetRequest(&Request{
			uri:    "/operations-manager/automations",
			params: &QueryParams{Limit: limit, Skip: skip},
		}, &res); err != nil {
			return nil, err
		}

		for _, ele := range res.Data {
			automations = append(automations, ele)
		}

		if len(automations) == res.Metadata.Total {
			break
		}

		skip += limit
	}

	logger.Info("Found %v automations", len(automations))

	return automations, nil
}

// Clear deletes all automations from the server
func (svc *AutomationService) Clear() error {
	logger.Trace()

	automations, err := svc.GetAll()
	if err != nil {
		return err
	}

	for _, ele := range automations {
		if err := svc.Delete(ele.Id); err != nil {
			return err
		}
	}

	return nil
}

// Import implements the `PUT /operations-manager/automations` requuest
func (svc *AutomationService) Import(in Automation) (*Automation, error) {
	logger.Trace()

	if len(in.Gbac.Read) > 0 && len(in.Gbac.Write) == 0 {
		return nil, errors.New("write group must be configured, when read group present")
	}

	var automations []any

	if len(in.Triggers) == 0 {
		b, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}

		var item map[string]interface{}
		if err := json.Unmarshal(b, &item); err != nil {
			return nil, err
		}

		item["triggers"] = []any{}

		automations = append(automations, item)
	} else {
		automations = append(automations, in)
	}

	body := map[string][]any{
		"automations": automations,
	}

	type Data struct {
		Success bool       `json:"success"`
		Data    Automation `json:"data"`
	}

	type Response struct {
		Data     []Data   `json:"data"`
		Message  string   `json:"message"`
		Metadata Metadata `json:"metadata"`
	}

	var res Response

	if err := svc.PutRequest(&Request{
		uri:  "/operations-manager/automations",
		body: &body,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	return &res.Data[0].Data, nil
}

// Export exports an automation by ID, including its triggers
func (svc *AutomationService) Export(id string) (*Automation, error) {
	logger.Trace()

	type Response struct {
		Data     *Automation `json:"data"`
		Message  string      `json:"message"`
		Metadata Metadata    `json:"metadata"`
	}

	var res Response
	var uri = fmt.Sprintf("/operations-manager/automations/%s/export", id)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	triggers := res.Data.Triggers
	res.Data.Triggers = []Trigger{}

	for _, ele := range triggers {
		var trigger Trigger

		b, err := json.Marshal(ele.(map[string]interface{}))
		if err != nil {
			logger.Fatal(err, "error trying to marshal data")
		}

		switch ele.(map[string]interface{})["type"].(string) {
		case "endpoint":
			var t EndpointTrigger
			if err := json.Unmarshal(b, &t); err != nil {
				logger.Fatal(err, "error trying to decode data")
			}
			trigger = t
		case "eventSystem":
			var t EventTrigger
			if err := json.Unmarshal(b, &t); err != nil {
				logger.Fatal(err, "error trying to decode data")
			}
			trigger = t
		case "manual":
			var t ManualTrigger
			if err := json.Unmarshal(b, &t); err != nil {
				logger.Fatal(err, "error trying to decode data")
			}
			trigger = t
		case "schedule":
			var t ScheduleTrigger
			if err := json.Unmarshal(b, &t); err != nil {
				logger.Fatal(err, "error trying to decode data")
			}
			trigger = t
		}

		res.Data.Triggers = append(res.Data.Triggers, trigger)
	}

	return res.Data, nil
}
