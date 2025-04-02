// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type UserRole struct {
	RoleId string `json:"roleId"`
}

type UserSettings struct {
	Id            string                 `json:"_id"`
	Meta          map[string]interface{} `json:"_meta"`
	AssignedRoles []UserRole             `json:"assignedRoles"`
	FirstName     string                 `json:"firstname"`
	GitTokens     map[string]interface{} `json:"gitTokens"`
	Inactive      bool                   `json:"inactive"`
	LastLogin     string                 `json:"lastLogin"`
	MemberOf      []interface{}          `json:"memberOf"`
	Provenance    string                 `json:"provenance"`
	Username      string                 `json:"username"`
}

type UserSettingsService struct {
	client *ServiceClient
}

func NewUserSettingsService(c client.Client) *UserSettingsService {
	return &UserSettingsService{
		client: NewServiceClient(c),
	}
}

func (svc *UserSettingsService) Get() (*UserSettings, error) {
	logger.Trace()

	var res *UserSettings

	if err := svc.client.Get("/user/settings", &res); err != nil {
		return nil, err
	}

	return res, nil
}
