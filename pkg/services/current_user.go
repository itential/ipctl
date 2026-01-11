// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
)

type UserGroup struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Provenance string `json:"provenance"`
}

type CurrentUser struct {
	Id         string                 `json:"id"`
	Username   string                 `json:"username"`
	FirstName  string                 `json:"firstname"`
	Groups     []UserGroup            `json:"groups"`
	Methods    map[string]interface{} `json:"methods"`
	Provenance string                 `json:"provenance"`
	Roles      []string               `json:"roles"`
	Routes     []string               `json:"routes"`
}

func GetCurrentUser(c client.Client) (*CurrentUser, error) {
	logging.Trace()

	svc := NewBaseService(c)

	var res *CurrentUser

	if err := svc.Get("/whoami", &res); err != nil {
		return nil, err
	}

	return res, nil
}
