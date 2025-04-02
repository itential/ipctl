// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"errors"
	"fmt"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

// Account represents a configured account configured on Itential Platform that
// can access the system.
type Account struct {
	// The account ID is a unique identifier that is assigned by the server and
	// can be used to identify a specific account in Platform.
	Id string `json:"_id"`

	// The email address for the account
	Email string `json:"email"`

	// The first name of the user for this account
	FirstName string `json:"firstname"`

	// Returns whether or not the account is inactive
	Inactive bool `json:"inactive"`

	// Returns whether or not the account is currently logged at the time of
	// the API call
	LoggedIn bool `json:"loggedIn"`

	// Identifies the origin of the account.  Since all accounts are federated
	// from other systems, this field provides an indication as to which
	// external system the account is sourced from.
	Provenance string `json:"provenance"`

	// The username associated with this account.  The username should be
	// unique within the system

	// The username associated with this account.  The username should be
	// unique within the system
	Username string `json:"username"`
}

// AccountService provides API access for mananging Itential Platform accounts.
type AccountService struct {
	client *ServiceClient
}

// Returns a new instance of AccountService using the conneciton as specified
// by c
func NewAccountService(c client.Client) *AccountService {
	return &AccountService{
		client: NewServiceClient(c),
	}
}

// GetAll calls `GET /authorization/accounts`
func (svc *AccountService) GetAll() ([]Account, error) {
	logger.Trace()

	type Response struct {
		Results []Account `json:"results"`
		Total   int       `json:"total"`
	}

	var res *Response

	if err := svc.client.Get("/authorization/accounts", &res); err != nil {
		return nil, err
	}

	return res.Results, nil
}

// Get invokes `GET /authorization/accounts/{id}`
func (svc *AccountService) Get(id string) (*Account, error) {
	logger.Trace()

	var res *Account

	var uri = fmt.Sprintf("/authorization/accounts/%s", id)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetByName accepts a single argument which is the username of the acccount to
// retrieve.  This function will iterate through all accounts and attempt to
// match the username against the returned data.  If a match is found, an
// Account is returned.  If a match is not found, an error is returned.
func (svc *AccountService) GetByName(name string) (*Account, error) {
	logger.Trace()

	accounts, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	var res *Account

	for _, ele := range accounts {
		if ele.Username == name {
			res = &ele
			break
		}
	}

	if res == nil {
		return nil, errors.New("account not found")
	}

	return res, nil
}

func (svc *AccountService) Deactivate(id string) error {
	logger.Trace()
	return svc.client.PatchRequest(&Request{
		uri:  fmt.Sprintf("/authorization/accounts/%s", id),
		body: map[string]interface{}{"inactive": true},
	}, nil)
}

func (svc *AccountService) Activate(id string) error {
	logger.Trace()
	return svc.client.PatchRequest(&Request{
		uri:  fmt.Sprintf("/authorization/accounts/%s", id),
		body: map[string]interface{}{"inactive": false},
	}, nil)
}
