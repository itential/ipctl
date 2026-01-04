// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type View struct {
	Provenance string `json:"provenance"`
	Path       string `json:"path"`
}

type ViewService struct {
	BaseService
}

func NewViewService(c client.Client) *ViewService {
	return &ViewService{BaseService: NewBaseService(c)}
}

// GetAll will retrieve all authorization views from the service and return
// them as an array.  If there are no configured authorization views, this
// function will return an empty array
func (svc *ViewService) GetAll() ([]View, error) {
	logger.Trace()

	type Response struct {
		Results []View `json:"results"`
		Total   int    `json:"total"`
	}

	var res *Response
	var views []View

	var limit = 100
	var skip = 0

	for {
		if err := svc.GetRequest(&Request{
			uri:    "/authorization/views",
			params: &QueryParams{Limit: limit, Skip: skip},
		}, &res); err != nil {
			return nil, err
		}

		for _, ele := range res.Results {
			views = append(views, ele)
		}

		if len(views) == res.Total {
			break
		}

		skip += limit
	}

	logger.Info("Found %v view(s)", len(views))

	return views, nil
}
