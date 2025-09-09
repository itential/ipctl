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

// RoleMethod represents an allowed method for a role
type RoleMethod struct {
	Name       string `json:"name"`
	Provenance string `json:"provenance"`
}

// RoleView represents an allowed view for a role
type RoleView struct {
	Provenance string `json:"provenance"`
	Path       string `json:"path"`
}

// Role represents a role in the authorization system
type Role struct {
	Id             string       `json:"_id,omitempty"`
	Provenance     string       `json:"provenance"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	AllowedMethods []RoleMethod `json:"allowedMethods"`
	AllowedViews   []RoleView   `json:"allowedViews"`
}

// RoleService provides methods for managing roles
type RoleService struct {
	client client.Client
}

// NewRoleService creates a new RoleService with the given client
func NewRoleService(c client.Client) *RoleService {
	return &RoleService{client: c}
}

// Create implements `http.MethodPost /authorization/roles`
func (svc *RoleService) Create(in Role) (*Role, error) {
	logger.Trace()

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    *Role  `json:"data"`
	}

	var response Response

	// NOTE (privateip): The description field must be in the body or the
	// server will return a  500 errror even though the description field is not
	// used.
	body := map[string]interface{}{
		"name":        in.Name,
		"provenance":  in.Provenance,
		"description": "",
	}

	if len(in.AllowedMethods) == 0 {
		body["allowedMethods"] = []any{}
	} else {
		body["allowedMethods"] = in.AllowedMethods
	}

	if len(in.AllowedViews) == 0 {
		body["allowedViews"] = []any{}
	} else {
		body["allowedViews"] = in.AllowedViews
	}

	resp, err := Do(&Request{
		client:             svc.client,
		method:             http.MethodPost,
		uri:                "/authorization/roles",
		body:               map[string]interface{}{"role": body},
		expectedStatusCode: http.StatusOK,
	})
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, err
	}

	logger.Info(response.Message)

	return response.Data, nil
}

// Delete implements `http.MethodDelete /authorization/roles/{id}`
func (svc *RoleService) Delete(id string) error {
	logger.Trace()

	resp, err := Do(&Request{
		client: svc.client,
		method: http.MethodDelete,
		uri:    fmt.Sprintf("/authorization/roles/%s", id),
	})

	if err != nil {
		return err
	}

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var response Response

	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return err
	}

	logger.Info(response.Message)

	return nil
}

// GetAll implements `http.MethodGet /authorization/roles`
func (svc *RoleService) GetAll() ([]Role, error) {
	logger.Trace()

	type Response struct {
		Results []Role `json:"results"`
		Total   int    `json:"total"`
	}

	var response *Response
	var roles []Role

	var limit = 100
	var skip = 0

	for {
		resp, err := Do(&Request{
			client:   svc.client,
			method:   http.MethodGet,
			uri:      "/authorization/roles",
			params:   &QueryParams{Limit: limit, Skip: skip},
			response: &response,
		})

		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(string(resp.Body))
		}

		for _, ele := range response.Results {
			roles = append(roles, ele)
		}

		if len(roles) == response.Total {
			break
		}

		skip += limit
	}

	logger.Info("Found %v roles", len(roles))

	return roles, nil
}

// Get implements `http.MethodGet /authorization/roles/{id}`
func (svc *RoleService) Get(id string) (*Role, error) {
	logger.Trace()

	var response map[string]interface{}

	res, err := Do(&Request{
		client:   svc.client,
		method:   http.MethodGet,
		uri:      fmt.Sprintf("/authorization/roles/%s", id),
		response: &response,
	})

	if err != nil {
		return nil, err
	}

	var role Role

	var body map[string]interface{}
	if err := json.Unmarshal(res.Body, &body); err != nil {
		logger.Fatal(err, "failed to unmarshal body")
	}

	if err := Unmarshal(body, &role); err != nil {
		return nil, err
	}

	return &role, nil

}

// Import imports a role into the authorization system
func (svc *RoleService) Import(in Role) (*Role, error) {
	logger.Trace()

	body := map[string]interface{}{
		"role": in,
	}

	res, err := Do(&Request{
		client:             svc.client,
		method:             http.MethodPost,
		uri:                "/authorization/roles",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	})

	if err != nil {
		return nil, err
	}

	type Response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Data    *Role  `json:"data"`
	}

	var response Response

	if err := json.Unmarshal(res.Body, &response); err != nil {
		return nil, err
	}

	logger.Info(response.Message)

	return response.Data, nil
}
