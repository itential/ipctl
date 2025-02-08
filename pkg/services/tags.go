// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"errors"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type Tag struct {
	Id          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TagService struct {
	client *ServiceClient
}

func NewTag(name, desc string) Tag {
	logger.Trace()
	return Tag{Name: name, Description: desc}

}

func NewTagService(iapClient client.Client) *TagService {
	return &TagService{client: NewServiceClient(iapClient)}
}

// GetAll implements `GET /tags/all`
func (svc *TagService) GetAll() ([]Tag, error) {
	logger.Trace()

	var res []Tag

	if err := svc.client.Get("/tags/all", &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *TagService) Get(id string) (*Tag, error) {
	logger.Trace()

	var res *Tag

	body := map[string]string{
		"id": id,
	}

	if err := svc.client.PostRequest(&Request{
		uri:                "/tags/get",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *TagService) GetByName(name string) (*Tag, error) {
	logger.Trace()

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
	logger.Trace()

	body := map[string]interface{}{
		"data": map[string]string{"ref_id": id},
	}

	var res []Tag

	if err := svc.client.PostRequest(&Request{
		uri:                "/tags/getTagsByReference",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *TagService) Create(in Tag) (*Tag, error) {
	logger.Trace()

	body := map[string]interface{}{
		"name":        in.Name,
		"description": in.Description,
	}

	var res *Tag

	if err := svc.client.PostRequest(&Request{
		uri:                "/tags/create",
		body:               map[string]interface{}{"data": body},
		expectedStatusCode: 200,
	}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *TagService) Delete(id string) error {
	logger.Trace()
	return svc.client.PostRequest(&Request{
		uri:                "/tags/delete",
		body:               map[string]interface{}{"_id": id},
		expectedStatusCode: 200,
	}, nil)

}
