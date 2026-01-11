// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"errors"
	"net/http"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
)

type Tag struct {
	Id          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TagService struct {
	BaseService
}

func NewTag(name, desc string) Tag {
	logging.Trace()
	return Tag{Name: name, Description: desc}

}

func NewTagService(c client.Client) *TagService {
	return &TagService{BaseService: NewBaseService(c)}
}

// GetAll implements `GET /tags/all`
func (svc *TagService) GetAll() ([]Tag, error) {
	logging.Trace()

	var res []Tag

	if err := svc.BaseService.Get("/tags/all", &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *TagService) Get(id string) (*Tag, error) {
	logging.Trace()

	var res *Tag

	body := map[string]string{
		"id": id,
	}

	if err := svc.PostRequest(&Request{
		uri:                "/tags/get",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *TagService) GetByName(name string) (*Tag, error) {
	logging.Trace()

	tags, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	var tag *Tag

	for _, ele := range tags {
		if ele.Name == name {
			tag = &ele
			break
		}
	}

	if tag == nil {
		return nil, errors.New("tag not found")
	}

	return svc.Get(tag.Id)
}

func (svc *TagService) GetTagsForReference(id string) ([]Tag, error) {
	logging.Trace()

	body := map[string]interface{}{
		"data": map[string]string{"ref_id": id},
	}

	var res []Tag

	if err := svc.PostRequest(&Request{
		uri:                "/tags/getTagsByReference",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *TagService) Create(in Tag) (*Tag, error) {
	logging.Trace()

	body := map[string]interface{}{
		"name":        in.Name,
		"description": in.Description,
	}

	var res *Tag

	if err := svc.PostRequest(&Request{
		uri:                "/tags/create",
		body:               map[string]interface{}{"data": body},
		expectedStatusCode: 200,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *TagService) Delete(id string) error {
	logging.Trace()
	return svc.PostRequest(&Request{
		uri:                "/tags/delete",
		body:               map[string]interface{}{"_id": id},
		expectedStatusCode: 200,
	}, nil)

}
