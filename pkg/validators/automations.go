// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package validators

import (
	"errors"
	"fmt"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type AutomationValidator struct {
	client  client.Client
	service *services.AutomationService
}

func NewAutomationValidator(c client.Client) AutomationValidator {
	return AutomationValidator{
		client:  c,
		service: services.NewAutomationService(c),
	}
}

func (v AutomationValidator) CanImport(in services.Automation) error {
	logger.Trace()

	if v.Exists(in.Name) {
		return errors.New(
			fmt.Sprintf("automation already exists"),
		)
	}

	if NewWorkflowValidator(v.client).Exists(in.ComponentName) {
		return errors.New(
			fmt.Sprintf("workflow `%s` does not exist on the destination server", in.ComponentName),
		)
	}

	if err := v.validateAccountsExist(in); err != nil {
		return err
	}

	return nil

}

func (v AutomationValidator) Exists(name string) bool {
	logger.Trace()

	exists, err := v.service.GetByName(name)
	if err != nil {
		if err.Error() != "automation not found" {
			logger.Fatal(err, "")
		}
	}
	return exists != nil
}

func (v AutomationValidator) validateAccountsExist(in services.Automation) error {
	logger.Trace()

	if err := v.checkAccounts(in.Gbac.Read); err != nil {
		return errors.New(
			fmt.Sprintf("account `%s` is assigned read permissions but does not exist on the destination server", err.(MissingAccountError).Username),
		)
	}

	if err := v.checkAccounts(in.Gbac.Write); err != nil {
		return errors.New(
			fmt.Sprintf("account `%s` is assigned write permissions but does not exist on the destination server", err.(MissingAccountError).Username),
		)
	}

	return nil
}

func (v AutomationValidator) checkAccounts(in []interface{}) error {
	var accounts []string

	for _, ele := range in {
		accounts = append(accounts, ele.(map[string]interface{})["name"].(string))
	}

	return NewAccountValidator(v.client).Exists(accounts)
}
