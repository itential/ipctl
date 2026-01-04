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

type DeviceGroup struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Devices       []string               `json:"devices"`
	Created       string                 `json:"created"`
	CreatedBy     string                 `json:"createdBy"`
	Updated       string                 `json:"updated"`
	LastUpdatedBy string                 `json:"lastUpdatedBy"`
	Gbac          map[string]interface{} `json:"gbac"`
}

type DeviceGroupService struct {
	BaseService
}

func NewDeviceGroup(name, desc string) DeviceGroup {
	logger.Trace()
	return DeviceGroup{
		Name:        name,
		Description: desc,
	}
}

func NewDeviceGroupService(c client.Client) *DeviceGroupService {
	return &DeviceGroupService{BaseService: NewBaseService(c)}
}

func (svc *DeviceGroupService) Get(id string) (*DeviceGroup, error) {
	logger.Trace()

	var res *DeviceGroup
	var uri = fmt.Sprintf("/configuration_manager/devicegroups/%s", id)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *DeviceGroupService) GetAll() ([]DeviceGroup, error) {
	logger.Trace()

	var res []DeviceGroup

	if err := svc.BaseService.Get("/configuration_manager/deviceGroups", &res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetByName retrieves a device group by name using client-side filtering.
// DEPRECATED: Business logic method - prefer using resources.DeviceGroupResource.GetByName
func (svc *DeviceGroupService) GetByName(name string) (*DeviceGroup, error) {
	logger.Trace()

	groups, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	for i := range groups {
		if groups[i].Name == name {
			return &groups[i], nil
		}
	}

	return nil, errors.New("device group not found")
}

func (svc *DeviceGroupService) Create(in DeviceGroup) (*DeviceGroup, error) {
	logger.Trace()

	body := map[string]interface{}{
		"groupName":        in.Name,
		"groupDescription": in.Description,
		"deviceNames":      "",
	}

	type Response struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var res Response

	if err := svc.PostRequest(&Request{
		uri:                "/configuration_manager/devicegroup",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info("%s", res.Message)

	return svc.Get(res.Id)
}

func (svc *DeviceGroupService) Delete(id string) error {
	logger.Trace()

	body := map[string]interface{}{
		"groupIds": []string{id},
	}

	type Response struct {
		Status  string `json:"status"`
		Deleted int    `json:"deleted"`
	}

	var res Response

	return svc.DeleteRequest(&Request{
		uri:                "/configuration_manager/devicegroups",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res)
}
