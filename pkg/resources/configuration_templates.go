// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"fmt"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/services"
)

// ConfigurationTemplateResource provides business logic for configuration template operations.
type ConfigurationTemplateResource struct {
	BaseResource
	service services.ConfigurationTemplateServicer
}

// NewConfigurationTemplateResource creates a new ConfigurationTemplateResource with the given service.
func NewConfigurationTemplateResource(svc services.ConfigurationTemplateServicer) ConfigurationTemplateResourcer {
	return &ConfigurationTemplateResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// GetByName retrieves a configuration template by name using client-side filtering.
// It fetches all templates and searches for a matching name.
func (r *ConfigurationTemplateResource) GetByName(name string) (*services.ConfigurationTemplate, error) {
	logging.Trace()

	templates, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	// Use service.Get to fetch full template details after finding by name
	for _, template := range templates {
		if template.Name == name {
			return r.service.Get(template.Id)
		}
	}

	logging.Error(nil, "configuration template not found")
	return nil, fmt.Errorf("configuration template not found")
}
