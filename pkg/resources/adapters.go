// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

// AdapterResource provides business logic for adapter operations.
type AdapterResource struct {
	BaseResource
	service services.AdapterServicer
}

// NewAdapterResource creates a new AdapterResource with the given service.
func NewAdapterResource(svc services.AdapterServicer) AdapterResourcer {
	return &AdapterResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// Restart orchestrates stopping and then starting an adapter.
// This is a composite operation that ensures proper adapter restart.
func (r *AdapterResource) Restart(name string) error {
	logger.Trace()

	if err := r.service.Stop(name); err != nil {
		return err
	}

	return r.service.Start(name)
}

// GetAll retrieves all configured adapter instances from the API.
// This is a pass-through to the service layer for pure API access.
func (r *AdapterResource) GetAll() ([]services.Adapter, error) {
	return r.service.GetAll()
}

// Get retrieves a specific adapter by name from the API.
// This is a pass-through to the service layer for pure API access.
func (r *AdapterResource) Get(name string) (*services.Adapter, error) {
	return r.service.Get(name)
}

// Create creates a new adapter instance with the provided configuration.
// This is a pass-through to the service layer for pure API access.
func (r *AdapterResource) Create(in services.Adapter) (*services.Adapter, error) {
	return r.service.Create(in)
}

// Delete removes the adapter instance with the specified name.
// This is a pass-through to the service layer for pure API access.
func (r *AdapterResource) Delete(name string) error {
	return r.service.Delete(name)
}

// Update modifies an existing adapter instance with the provided configuration.
// This is a pass-through to the service layer for pure API access.
func (r *AdapterResource) Update(in services.Adapter) (*services.Adapter, error) {
	return r.service.Update(in)
}

// Export retrieves the adapter configuration for backup or import operations.
// This is a pass-through to the service layer for pure API access.
func (r *AdapterResource) Export(name string) (*services.Adapter, error) {
	return r.service.Export(name)
}

// Start initiates the specified adapter instance.
// This is a pass-through to the service layer for pure API access.
func (r *AdapterResource) Start(name string) error {
	return r.service.Start(name)
}

// Stop halts the specified adapter instance.
// This is a pass-through to the service layer for pure API access.
func (r *AdapterResource) Stop(name string) error {
	return r.service.Stop(name)
}

// Import imports an adapter configuration.
// This is a pass-through to the service layer for pure API access.
func (r *AdapterResource) Import(in services.Adapter) (*services.Adapter, error) {
	return r.service.Import(in)
}
