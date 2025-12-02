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

type GoldenConfigTreeSummary struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	DeviceType  string                 `json:"deviceType"`
	Versions    []string               `json:"versions"`
	Gbac        map[string]interface{} `json:"gbac"`
	Created     string                 `json:"created"`
	LastUpdated string                 `json:"lastUpdated"`
	Tags        []Tag                  `json:"tags"`
}

type GoldenConfigTree struct {
	Id            string                 `json:"_id"`
	Name          string                 `json:"name"`
	TreeId        string                 `json:"treeId"`
	Version       string                 `json:"version"`
	DeviceType    string                 `json:"deviceType"`
	Root          map[string]interface{} `json:"root"`
	Created       string                 `json:"created"`
	CreatedBy     string                 `json:"createdBy"`
	LastUpdated   string                 `json:"lastUpdated"`
	LastUpdatedBy string                 `json:"lastUpdatedBy"`
	Variables     map[string]interface{} `json:"variables"`
	Base          string                 `json:"base,omitempty"`
	Tags          []Tag                  `json:"tags"`
}

type GoldenConfigService struct {
	client *ServiceClient
}

func NewGoldenConfigService(c client.Client) *GoldenConfigService {
	return &GoldenConfigService{client: NewServiceClient(c)}
}

// Create calls `POST /configuration_manager/configs`
func (svc *GoldenConfigService) Create(in GoldenConfigTree) (*GoldenConfigTree, error) {
	logger.Trace()

	body := map[string]interface{}{
		"name":       in.Name,
		"deviceType": in.DeviceType,
	}

	var res *GoldenConfigTree

	if err := svc.client.PostRequest(&Request{
		uri:                "/configuration_manager/configs",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// Delete calls `DELETE /configuration_manager/configs/{id}`
func (svc *GoldenConfigService) Delete(id string) error {
	logger.Trace()
	var uri = fmt.Sprintf("/configuration_manager/configs/%s", id)
	return svc.client.Delete(uri)
}

// GetAll calls `GET /configuration_manager/configs`
func (svc *GoldenConfigService) GetAll() ([]GoldenConfigTreeSummary, error) {
	logger.Trace()

	var res []GoldenConfigTreeSummary
	var uri = "/configuration_manager/configs"

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *GoldenConfigService) GetByName(name string) (*GoldenConfigTreeSummary, error) {
	logger.Trace()

	trees, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	var res *GoldenConfigTreeSummary

	for _, ele := range trees {
		if ele.Name == name {
			res = &ele
			break
		}
	}

	if res == nil {
		return nil, errors.New("gctree not found")
	}

	return res, nil
}

// Import will attempt to import a golden configuraiton tree specified by the
// `in` argument and import it to the server.  This function will return an
// error or nil if no error is encountered
func (svc *GoldenConfigService) Import(in GoldenConfigTree) error {
	logger.Trace()

	body := map[string]interface{}{
		"trees": []map[string]interface{}{
			map[string]interface{}{
				"data": []GoldenConfigTree{in},
			},
		},
		"options": map[string]interface{}{},
	}

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/configuration_manager/import/goldenconfigs",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return err
	}

	logger.Info("%s", res.Message)

	return nil
}

// Export calls `POST /configuration_manager/export/goldenconfigs`
func (svc *GoldenConfigService) Export(id string) (*GoldenConfigTree, error) {
	logger.Trace()

	body := map[string]interface{}{"treeId": id}

	type Response struct {
		Status string             `json:"status"`
		Data   []GoldenConfigTree `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/configuration_manager/export/goldenconfigs",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return &res.Data[0], nil
}
