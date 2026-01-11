// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package validators

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/services"
)

type WorkflowValidator struct {
	client  client.Client
	service *services.WorkflowService
}

func NewWorkflowValidator(c client.Client) WorkflowValidator {
	return WorkflowValidator{
		client:  c,
		service: services.NewWorkflowService(c),
	}
}

func (v WorkflowValidator) Exists(name string) bool {
	logging.Trace()

	res, err := v.service.Get(name)
	if err != nil {
		if err.Error() == "workflow not found" {
			return false
		} else {
			logging.Fatal(err, "")
		}
	}

	return res == nil
}
