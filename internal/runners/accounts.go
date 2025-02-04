// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type AccountRunner struct {
	config  *config.Config
	service *services.AccountService
}

func NewAccountRunner(client client.Client, cfg *config.Config) *AccountRunner {
	return &AccountRunner{
		config:  cfg,
		service: services.NewAccountService(client),
	}
}

// Get implements the `get accounts` command
func (r *AccountRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	accounts, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tSOURCE"}
	for _, ele := range accounts {
		lines := []string{ele.Username, ele.Provenance}
		display = append(display, strings.Join(lines, "\t"))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(accounts),
	), nil
}

// Describe implements the `describe account <name>` command
func (r *AccountRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	accounts, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var account *services.Account

	for _, ele := range accounts {
		if ele.Username == name {
			account = &ele
			break
		}
	}

	if account == nil {
		return nil, errors.New(fmt.Sprintf("account `%s` does not exist", name))
	}

	return NewResponse(
		"",
		WithJson(account),
	), nil
}
