// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"fmt"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/resources"
	"github.com/itential/ipctl/pkg/services"
)

type AccountRunner struct {
	BaseRunner
	resource resources.AccountResourcer
}

func NewAccountRunner(client client.Client, cfg *config.Config) *AccountRunner {
	return &AccountRunner{
		BaseRunner: NewBaseRunner(client, cfg),
		resource:   resources.NewAccountResource(services.NewAccountService(client)),
	}
}

/*
i******************************************************************************
Reader interface
*******************************************************************************
*/

// Get implements the `get accounts` command
func (r *AccountRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	accounts, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"username", "provenance"},
		Object: accounts,
	}, nil
}

// Describe implements the `describe account <name>` command
func (r *AccountRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	account, err := r.resource.GetByName(in.Args[0])
	if err != nil {
		return nil, fmt.Errorf("account `%s` does not exist", in.Args[0])
	}

	tmpl, err := templates.ReadFile("templates/accounts/describe.tmpl")
	if err != nil {
		logger.Fatal(err, "failed to load template")
	}

	return &Response{
		Object:   account,
		Template: string(tmpl),
	}, nil
}
