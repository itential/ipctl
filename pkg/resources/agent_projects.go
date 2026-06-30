// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/services"
)

// AgentProjectResource provides business logic for agent project operations.
type AgentProjectResource struct {
	BaseResource
	service services.AgentProjectServicer
}

// NewAgentProjectResource creates a new AgentProjectResource with the given service.
func NewAgentProjectResource(svc services.AgentProjectServicer) AgentProjectResourcer {
	return &AgentProjectResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// GetAll retrieves all agent projects.
func (r *AgentProjectResource) GetAll() ([]services.AgentProject, error) {
	return r.service.GetAll()
}

// Get retrieves an agent project by ID.
func (r *AgentProjectResource) Get(id string) (*services.AgentProject, error) {
	return r.service.Get(id)
}

// GetByName retrieves an agent project by name.
func (r *AgentProjectResource) GetByName(name string) (*services.AgentProject, error) {
	logging.Trace()
	return r.service.GetByName(name)
}

// Export exports an agent project bundle by project ID.
func (r *AgentProjectResource) Export(id string) (*services.AgentProjectBundle, error) {
	return r.service.Export(id)
}

// Import imports an agent project bundle.
func (r *AgentProjectResource) Import(bundle services.AgentProjectBundle) (*services.AgentProjectBundle, error) {
	return r.service.Import(bundle)
}
