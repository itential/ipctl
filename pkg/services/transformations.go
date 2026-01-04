// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type TransformationView struct {
	Col int `json:"col"`
	Row int `json:"row"`
}

type Transformation struct {
	Id          string                   `json:"_id,omitempty"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Functions   []map[string]interface{} `json:"functions"`
	Incoming    []map[string]interface{} `json:"incoming"`
	Outgoing    []map[string]interface{} `json:"outgoing"`
	Steps       []map[string]interface{} `json:"steps"`
	View        TransformationView       `json:"view"`
	Created     string                   `json:"created,omitempty"`
	LastUpdated string                   `json:"lastUpdated,omitempty"`
	Version     string                   `json:"version,omitempty"`
	Tags        []string                 `json:"tags"`
}

type TransformationService struct {
	BaseService
}

func NewTransformation(name, description string) Transformation {
	logger.Trace()

	return Transformation{
		Name:        name,
		Description: description,
		Functions:   []map[string]interface{}{},
		Incoming:    []map[string]interface{}{},
		Outgoing:    []map[string]interface{}{},
		Steps:       []map[string]interface{}{},
		Tags:        []string{},
		View:        TransformationView{Col: 3, Row: 5},
	}
}

func NewTransformationService(c client.Client) *TransformationService {
	return &TransformationService{BaseService: NewBaseService(c)}
}

// GetAll will retrieve all configured transformations from the server and
// return them as an array.  If there are no configured transformations, this
// function will return an empty array.
func (svc *TransformationService) GetAll() ([]Transformation, error) {
	logger.Trace()

	type Response struct {
		Results []Transformation `json:"results"`
		Total   int              `json:"total"`
	}

	var res Response
	var transformations []Transformation

	var limit = 100
	var skip = 0

	for {
		if err := svc.GetRequest(&Request{
			uri:    "/transformations",
			params: &QueryParams{Limit: limit, Skip: skip},
		}, &res); err != nil {
			return nil, err
		}

		for _, ele := range res.Results {
			transformations = append(transformations, ele)
		}

		if len(transformations) == res.Total {
			break
		}

		skip += limit
	}

	logger.Info("Found%v transformations", len(transformations))

	return transformations, nil
}

// Get will retrieve the transformation from the server with the specifid
// identifier.  If the transformation does not exist, this function will return
// a "transformation not found" error.
func (svc *TransformationService) Get(id string) (*Transformation, error) {
	logger.Trace()

	var res *Transformation
	var uri = fmt.Sprintf("/transformations/%s", id)

	if err := svc.BaseService.Get(uri, &res); err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("transformation not found")
	}

	return res, nil
}

// GetByName will attempt to find a transformation based on its name.  If the
// specified transformation is found, it will be returned to the calling
// function.  If a match is not found, this function will return a
// "transformation not found" error.
func (svc *TransformationService) GetByName(name string) (*Transformation, error) {
	logger.Trace()

	type Response struct {
		Results []Transformation `json:"results"`
		Total   int              `json:"total"`
	}

	var res Response

	if err := svc.GetRequest(&Request{
		uri: "/transformations",
		query: map[string]string{
			"contains[name]": name,
		},
	}, &res); err != nil {
		return nil, err
	}

	if res.Total == 0 {
		return nil, errors.New("transformation not found")
	}

	var selected *Transformation

	for _, ele := range res.Results {
		if !strings.HasPrefix(ele.Name, "@") {
			selected = &ele
			break
		}
	}

	if selected == nil {
		return nil, errors.New("transformation not found")
	}

	return selected, nil
}

// Create all attempt to create a new transformation on the server.  This
// function accepts a single argument which is the transformation document.
// If transformation of the same name already exists on the server, this
// function will return an error.
func (svc *TransformationService) Create(in Transformation) (*Transformation, error) {
	logger.Trace()

	var res *Transformation

	if err := svc.PostRequest(&Request{
		uri:                "/transformations",
		body:               &in,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// Delete will remove an existing transformation with the specifie ID from the
// server.  If the specified ID does not exist, this function will still return
// successfully.
func (svc *TransformationService) Delete(id string) error {
	logger.Trace()
	return svc.DeleteRequest(&Request{
		uri:                fmt.Sprintf("/transformations/%s", id),
		expectedStatusCode: http.StatusNoContent,
	}, nil)
}

func (svc *TransformationService) Clear() error {
	logger.Trace()

	elements, err := svc.GetAll()
	if err != nil {
		return err
	}

	for _, ele := range elements {
		if err := svc.Delete(ele.Id); err != nil {
			return err
		}
	}

	return nil
}

// Import will attempt to import a transformation into the current server.
// This function accepts a single argument of type Transformation.  If no
// transformation with the same name exists, the document will be imported.  If
// an exisitng tranformation with the same name exists, this function will
// successfully import the document but change the name to append an
// incremented value.  For instance if importing a transformation named "test"
// and the name is not unique, the server API will change the transformation
// name to "test (1)"
func (svc *TransformationService) Import(in Transformation) (*Transformation, error) {
	logger.Trace()

	var res *Transformation

	if err := svc.PostRequest(&Request{
		uri:                "/transformations/import",
		body:               &in,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}
