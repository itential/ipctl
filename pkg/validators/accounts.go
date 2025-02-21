// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package validators

import (
	"fmt"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type MissingAccountError struct {
	Username string
}

func (e MissingAccountError) Error() string {
	return fmt.Sprintf("account `%s` does not exist on the server", e.Username)
}

func NewMissingAccountError(u string) MissingAccountError {
	return MissingAccountError{Username: u}
}

type AccountValidator struct {
	client  client.Client
	service *services.AccountService
}

func NewAccountValidator(c client.Client) AccountValidator {
	return AccountValidator{
		client:  c,
		service: services.NewAccountService(c),
	}
}

func (v AccountValidator) Exists(accounts []string) error {
	logger.Trace()

	res, err := services.NewAccountService(v.client).GetAll()
	if err != nil {
		logger.Fatal(err, "")
	}

	var existing []string
	for _, ele := range res {
		existing = append(existing, ele.Username)
	}

	for _, ele := range accounts {
		if !StringInSlice(ele, existing) {
			return NewMissingAccountError(ele)
		}
	}

	return nil
}
