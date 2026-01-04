// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type MethodRoute struct {
	Path     string `json:"path"`
	FullPath string `json:"fullPath"`
	Verb     string `json:"verb"`
}

type Method struct {
	Provenance string      `json:"provenance"`
	Name       string      `json:"name"`
	Deprecated bool        `json:"deprecated"`
	Route      MethodRoute `json:"route"`
}

type MethodService struct {
	BaseService
}

func NewMethodService(c client.Client) *MethodService {
	return &MethodService{BaseService: NewBaseService(c)}
}

// GetAll implements `GET /authorization/methods`
func (svc *MethodService) GetAll() ([]Method, error) {
	logger.Trace()

	type Response struct {
		Results []Method `json:"results"`
		Total   int      `json:"total"`
	}

	var res *Response

	if err := svc.BaseService.Get("/authorization/methods", &res); err != nil {
		return nil, err
	}

	return res.Results, nil
}
