// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package validators

import (
	"errors"
	"fmt"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/services"
)

type ProfileValidator struct {
	client  client.Client
	service *services.ProfileService
}

func NewProfileValidator(c client.Client) ProfileValidator {
	return ProfileValidator{
		client:  c,
		service: services.NewProfileService(c),
	}
}

// CanImport will check the current server to see if the profile can be
// imported successfuly.   This function will check if a profile with the same
// name exists on the server.  If one does, this function will return a
// ValidationError.  If one doesn't, nil will be returned indiciating the
// profile can be successfully imported.
func (v ProfileValidator) CanImport(p services.Profile) *ValidatorError {
	logging.Trace()

	existing, err := v.service.Get(p.Id)
	if err != nil {
		return ServiceError(err)
	}

	if existing != nil {
		return AssetExists(
			errors.New(
				fmt.Sprintf("profile `%s` already exists on the server", p.Id),
			),
		)
	}

	return nil
}
