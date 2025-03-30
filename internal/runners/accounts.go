// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"embed"
	"fmt"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

//go:embed templates/accounts/*.tmpl
var content embed.FS

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

/*
i******************************************************************************
Reader interface
*******************************************************************************
*/

// Get implements the `get accounts` command
func (r *AccountRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	accounts, err := r.service.GetAll()
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

	accounts, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	account := r.selectAccountByName(in.Args[0], accounts)

	if account == nil {
		return nil, fmt.Errorf("account `%s` does not exist", in.Args[0])
	}

	tmpl, err := content.ReadFile("describe.tmpl")
	if err != nil {
		logger.Fatal(err, "failed to load template")
	}

	return &Response{
		Object:   account,
		Template: string(tmpl),
		//Template: "Username: {{.Username}}",
	}, nil
}

// selectAccountByUsername takes a list of service.Accounts and iterates over
// them looking for the first instance of username.   If found, the Account is
// returned.  If the username is not found, nil is returend.
func (r *AccountRunner) selectAccountByName(username string, accounts []services.Account) *services.Account {
	logger.Trace()
	for _, ele := range accounts {
		if ele.Username == username {
			return &ele
		}
	}
	return nil
}
