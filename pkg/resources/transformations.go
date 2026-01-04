// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"fmt"
	"strings"

	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

// TransformationResource provides business logic for transformation operations.
type TransformationResource struct {
	BaseResource
	service services.TransformationServicer
}

// NewTransformationResource creates a new TransformationResource with the given service.
func NewTransformationResource(svc services.TransformationServicer) TransformationResourcer {
	return &TransformationResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// GetAll retrieves all transformations from the API.
// This is a pass-through to the service layer for pure API access.
func (r *TransformationResource) GetAll() ([]services.Transformation, error) {
	return r.service.GetAll()
}

// Get retrieves a specific transformation by name from the API.
// This is a pass-through to the service layer for pure API access.
func (r *TransformationResource) Get(name string) (*services.Transformation, error) {
	return r.service.Get(name)
}

// Create creates a new transformation.
// This is a pass-through to the service layer for pure API access.
func (r *TransformationResource) Create(in services.Transformation) (*services.Transformation, error) {
	return r.service.Create(in)
}

// Delete removes a transformation by its identifier.
// This is a pass-through to the service layer for pure API access.
func (r *TransformationResource) Delete(id string) error {
	return r.service.Delete(id)
}

// Import imports a transformation.
// This is a pass-through to the service layer for pure API access.
func (r *TransformationResource) Import(in services.Transformation) (*services.Transformation, error) {
	return r.service.Import(in)
}

// GetByName retrieves a transformation by name using client-side filtering.
// It excludes system resources (those with names starting with "@").
// This method fetches all transformations and filters for the matching name
// while excluding system-managed transformations.
func (r *TransformationResource) GetByName(name string) (*services.Transformation, error) {
	logger.Trace()

	transformations, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	// Find transformation by name, excluding system resources
	for i := range transformations {
		if transformations[i].Name == name && !strings.HasPrefix(transformations[i].Name, "@") {
			return &transformations[i], nil
		}
	}

	logger.Error(nil, "transformation not found")
	return nil, fmt.Errorf("transformation not found")
}

// Clear deletes all transformations from the server.
// This is a bulk operation that orchestrates multiple delete calls.
func (r *TransformationResource) Clear() error {
	logger.Trace()

	transformations, err := r.service.GetAll()
	if err != nil {
		return err
	}

	return DeleteAll(transformations, func(t services.Transformation) string {
		return t.Id
	}, r.service.Delete)
}
