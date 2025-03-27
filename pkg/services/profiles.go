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

type ApplicationProperties struct {
	Directory string `json:"directory"`
}

type UIProperties struct {
	Description    string `json:"description"`
	Layout         string `json:"layout"`
	Home           string `json:"home"`
	Login          string `json:"login"`
	Profile        string `json:"profile"`
	UserConfig     string `json:"user_config"`
	GroupConfig    string `json:"group_config"`
	NewUser        string `json:"new_user"`
	EditUser       string `json:"edit_user"`
	NewGroup       string `json:"new_group"`
	FavIcon        string `json:"fav_icon"`
	AppleTouchIcon string `json:"apple_touch_icon"`
}

type ProfileMetadata struct {
	ActiveSync bool `json:"activeSync"`
	IsActive   bool `json:"isActive"`
}

type Profile struct {
	Id              string                   `json:"id"`
	Description     string                   `json:"description"`
	LaunchDelay     int                      `json:"launchDelay"`
	LaunchTimeout   int                      `json:"launchTimeout"`
	Application     ApplicationProperties    `json:"applicationProps"`
	UI              UIProperties             `json:"uiProps"`
	Authentication  map[string]interface{}   `json:"authenticationProps,omitempty"`
	Express         map[string]interface{}   `json:"expressProps,omitempty"`
	Logger          map[string]interface{}   `json:"loggerProps,omitempty"`
	Alarm           []map[string]interface{} `json:"alarmProps,omitempty"`
	Audit           map[string]interface{}   `json:"auditProps,omitempty"`
	Path            map[string]interface{}   `json:"pathProps,omitempty"`
	TaskWorker      map[string]interface{}   `json:"taskWorkerProps,omitempty"`
	Redis           map[string]interface{}   `json:"redisProps,omitempty"`
	Services        []string                 `json:"services,omitempty"`
	System          map[string]interface{}   `json:"systemProps,omitempty"`
	Integration     map[string]interface{}   `json:"integrationProps,omitempty"`
	Prebuilt        map[string]interface{}   `json:"prebuiltProps,omitempty"`
	RetryStrategy   map[string]interface{}   `json:"retryStrategy,omitempty"`
	AdapterStrategy map[string]interface{}   `json:"adapterStrategy,omitempty"`
	Adapter         map[string]interface{}   `json:"adapterProps,omitempty"`
}

type ProfileService struct {
	client *ServiceClient
}

func NewProfileService(iapClient client.Client) *ProfileService {
	return &ProfileService{client: NewServiceClient(iapClient)}
}

func NewProfile(name, desc string) Profile {
	logger.Trace()

	return Profile{
		Id:          name,
		Description: desc,
	}
}

func (svc *ProfileService) GetAll() ([]Profile, error) {
	logger.Trace()

	type Results struct {
		Metadata ProfileMetadata `json:"metadata"`
		Profile  Profile         `json:"profile"`
	}

	type Response struct {
		Results []Results `json:"results"`
		Total   int       `json:"total"`
	}

	var res Response

	if err := svc.client.Get("/profiles", &res); err != nil {
		return nil, err
	}

	var profiles []Profile
	for _, ele := range res.Results {
		profiles = append(profiles, ele.Profile)
	}

	return profiles, nil
}

func (svc *ProfileService) Get(id string) (*Profile, error) {
	logger.Trace()

	type Response struct {
		Metadata struct {
			IsActive   bool `json:"isActive"`
			ActiveSync bool `json:"activeSync"`
		} `json:"metadata"`
		Profile Profile `json:"profile"`
	}

	var res Response
	var uri = fmt.Sprintf("/profiles/%s", id)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return &res.Profile, nil
}

func (svc *ProfileService) GetActiveProfile() (*Profile, error) {
	logger.Trace()

	type Results struct {
		Metadata ProfileMetadata `json:"metadata"`
		Profile  Profile         `json:"profile"`
	}

	type Response struct {
		Results []Results `json:"results"`
		Total   int       `json:"total"`
	}

	var res Response

	if err := svc.client.Get("/profiles", &res); err != nil {
		return nil, err
	}

	var active *Profile

	for _, ele := range res.Results {
		if ele.Metadata.IsActive {
			active = &ele.Profile
			break
		}
	}

	if active == nil {
		return nil, errors.New("failed to find the active profile")
	}

	return active, nil
}

func (svc *ProfileService) Create(in Profile) (*Profile, error) {
	logger.Trace()

	properties := map[string]interface{}{
		"id":          in.Id,
		"description": in.Description,
	}

	type Response struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    *Profile `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/profiles",
		body:               map[string]interface{}{"properties": properties},
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil
}

func (svc *ProfileService) Delete(name string) error {
	logger.Trace()
	return svc.client.Delete(fmt.Sprintf("/profiles/%s", name))
}

func (svc *ProfileService) Import(in Profile) (*Profile, error) {
	logger.Trace()

	type Response struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    *Profile `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/profiles/import",
		body:               map[string]interface{}{"properties": in},
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (svc *ProfileService) Export(name string) (*Profile, error) {
	logger.Trace()

	type Response struct {
		Metadata Metadata `json:"metadata"`
		Profile  *Profile `json:"profile"`
	}

	var res Response

	// XXX (privateip): The export URI returns an error from the server so
	// using the get profile route instead
	//var uri = fmt.Sprintf("/profiles/%s/export", name)
	var uri = fmt.Sprintf("/profiles/%s", name)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res.Profile, nil
}

func (svc *ProfileService) Activate(name string) error {
	logger.Trace()

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var res Response
	var uri = fmt.Sprintf("/profiles/%s/active", name)

	if err := svc.client.Put(uri, nil, &res); err != nil {
		return err
	}

	return nil
}
