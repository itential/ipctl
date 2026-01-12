// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package resources provides the business logic layer for interacting with Itential Platform resources.
//
// This package sits between the runner layer and the service layer, implementing all business logic
// for resource operations. Resources consume pure API services and provide higher-level operations
// with business rules, validation, transformation, and orchestration.
//
// # Architecture
//
// Resources are part of a layered architecture:
//
//	Handler → Runner → Resource → Service → API
//
// Each resource wraps one or more services and adds business logic on top of pure API operations.
//
// # Interface-Based Design
//
// All resources implement corresponding interfaces (e.g., AccountResourcer, ProjectResourcer).
// This enables:
//   - Easy unit testing with mocked service dependencies
//   - Dependency injection
//   - Clear separation of concerns
//
// Example:
//
//	type AccountResourcer interface {
//	    GetAll() ([]services.Account, error)
//	    Get(id string) (*services.Account, error)
//	    GetByName(name string) (*services.Account, error)  // Business logic
//	    Activate(id string) error
//	    Deactivate(id string) error
//	}
//
// # Responsibilities
//
// Resources should implement:
//   - Client-side filtering and search (GetByName methods)
//   - Data transformation for import/export
//   - Business rule validation
//   - Multi-step orchestration (e.g., DeleteWithOptions)
//   - Complex operations spanning multiple services
//
// Resources should NOT:
//   - Make direct HTTP calls (delegate to services)
//   - Handle CLI-specific concerns (belongs in handlers/runners)
//   - Manage configuration (belongs in runners)
//   - Format output for display (belongs in runners)
//
// # Usage Examples
//
// Basic usage with service interface:
//
//	accountSvc := services.NewAccountService(client)
//	accountRes := resources.NewAccountResource(accountSvc)
//
//	// Pure API operation (pass-through)
//	accounts, err := accountRes.GetAll()
//
//	// Business logic operation (client-side filtering)
//	admin, err := accountRes.GetByName("admin")
//
// Complex orchestration:
//
//	modelSvc := services.NewModelService(client)
//	wfSvc := services.NewWorkflowService(client)
//	jstSvc := services.NewTransformationService(client)
//	instSvc := services.NewInstanceService(client)
//
//	modelRes := resources.NewModelResource(modelSvc, wfSvc, jstSvc, instSvc)
//
//	// Orchestrates multiple operations with business rules
//	err := modelRes.DeleteWithOptions(model, resources.DeleteOptions{
//	    DeleteInstances: true,
//	    DeleteRelated:   true,
//	})
//
// # Pass-Through vs Business Logic
//
// Resources provide two types of methods:
//
// 1. Pass-through methods - Direct delegation to service for pure API access:
//
//	func (r *AccountResource) GetAll() ([]services.Account, error) {
//	    return r.service.GetAll()
//	}
//
// 2. Business logic methods - Add validation, filtering, orchestration:
//
//	func (r *AccountResource) GetByName(name string) (*services.Account, error) {
//	    accounts, err := r.service.GetAll()
//	    if err != nil {
//	        return nil, err
//	    }
//	    // Client-side filtering logic
//	    return FindByName(accounts, name, func(a services.Account) string {
//	        return a.Username
//	    })
//	}
//
// # Common Patterns
//
// Client-side filtering (when API doesn't support it):
//
//	func (r *Resource) GetByName(name string) (*Type, error) {
//	    items, err := r.service.GetAll()
//	    // Apply filter logic
//	}
//
// Data transformation:
//
//	func (r *ProjectResource) Import(in services.Project) (*services.Project, error) {
//	    // Transform data structure
//	    transformed := r.transformImport(in)
//	    return r.service.Import(transformed)
//	}
//
// Multi-service orchestration:
//
//	func (r *ModelResource) DeleteWithOptions(model *services.Model, opts DeleteOptions) error {
//	    // Check instances
//	    // Delete related workflows
//	    // Delete related transformations
//	    // Finally delete model
//	}
//
// # Helper Functions
//
// The package provides generic helper functions in base.go:
//
//   - FindByName[T] - Generic search by name
//   - DeleteAll[T] - Generic bulk delete with error collection
//   - ValidateGbacRules - GBAC validation
//
// # Testing
//
// Resources are designed for easy testing. Create mock service interfaces:
//
//	type MockAccountService struct {
//	    mock.Mock
//	}
//
//	func (m *MockAccountService) GetAll() ([]services.Account, error) {
//	    args := m.Called()
//	    return args.Get(0).([]services.Account), args.Error(1)
//	}
//
//	// Test resource with mocked service
//	mockSvc := &MockAccountService{}
//	mockSvc.On("GetAll").Return([]services.Account{{Username: "test"}}, nil)
//	resource := resources.NewAccountResource(mockSvc)
//	account, err := resource.GetByName("test")
//
// # Error Handling
//
// Resources should wrap errors with context:
//
//	accounts, err := r.service.GetAll()
//	if err != nil {
//	    return nil, fmt.Errorf("resource layer: fetching accounts: %w", err)
//	}
package resources
