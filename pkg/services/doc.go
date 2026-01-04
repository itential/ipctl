// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package services provides pure API implementations for interacting with
// Itential Platform REST APIs.
//
// This package contains service implementations for all major Itential Platform
// resources including projects, workflows, automations, adapters, accounts, and more.
// Each service corresponds to a specific API endpoint and provides direct HTTP operations
// without business logic. Business logic resides in the resources package.
//
// # Architecture
//
// Services are part of a layered architecture:
//
//	Handler → Runner → Resource → Service → API
//
// All services embed the BaseService type which provides common HTTP operations:
//
//	type ProjectService struct {
//	    BaseService
//	}
//
// Services implement the Servicer interfaces defined in interfaces.go and use the
// client.Client interface for HTTP communication, returning structured domain objects
// specific to each resource type.
//
// # Interface-Based Design
//
// All services implement corresponding interfaces (e.g., ProjectServicer, AccountServicer).
// This enables:
//   - Easy unit testing with mocks
//   - Dependency injection
//   - Alternative implementations (caching, logging, etc.)
//
// Example:
//
//	type AccountServicer interface {
//	    GetAll() ([]Account, error)
//	    Get(id string) (*Account, error)
//	    Activate(id string) error
//	    Deactivate(id string) error
//	}
//
// # Separation of Concerns
//
// Services are pure API implementations and should:
//   - Make direct HTTP calls to the Itential Platform API
//   - Marshal/unmarshal JSON payloads
//   - Handle HTTP-level errors
//   - Return structured data types
//
// Services should NOT:
//   - Implement business logic (use resources package)
//   - Perform client-side filtering (use resources package)
//   - Validate business rules (use resources package)
//   - Orchestrate multiple API calls (use resources package)
//
// # Usage
//
// Create a service by passing a configured HTTP client:
//
//	client := client.NewHttpClient(profile)
//	projectSvc := services.NewProjectService(client)
//	projects, err := projectSvc.GetAll()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Error Handling
//
// Services return errors in the following cases:
//   - Network or connection failures
//   - HTTP status codes indicating errors (4xx, 5xx)
//   - JSON unmarshaling failures
//   - Resource not found conditions
//
// Errors can be checked for specific conditions:
//
//	project, err := projectSvc.GetByName("my-project")
//	if err != nil {
//	    if strings.Contains(err.Error(), "not found") {
//	        // Handle not found case
//	    }
//	    return err
//	}
//
// # Service Types
//
// The package provides services for the following resource categories:
//
// Automation Studio:
//   - ProjectService: Manage automation projects
//   - WorkflowService: Manage workflows
//   - TransformationService: Manage data transformations
//   - JsonFormService: Manage JSON forms
//   - TemplateService: Manage templates
//   - CommandTemplateService: Manage command templates
//   - AnalyticTemplateService: Manage analytic templates
//
// Operations Manager:
//   - AutomationService: Manage automations and orchestrations
//
// Admin Essentials:
//   - AccountService: Manage user accounts
//   - GroupService: Manage user groups
//   - RoleService: Manage roles and permissions
//   - ProfileService: Manage configuration profiles
//   - AdapterService: Manage adapter instances
//   - IntegrationService: Manage integrations
//
// Configuration Manager:
//   - DeviceService: Manage network devices
//   - DeviceGroupService: Manage device groups
//   - ConfigParserService: Manage configuration parsers
//   - GoldenConfigService: Manage golden configurations
//
// Lifecycle Manager:
//   - ModelService: Manage LCM models and instances
//
// # Common Patterns
//
// Most services implement similar patterns:
//
// GetAll - Retrieve all resources with automatic pagination:
//
//	resources, err := service.GetAll()
//
// Get - Retrieve a single resource by ID:
//
//	resource, err := service.Get(id)
//
// GetByName - Retrieve a resource by name:
//
//	resource, err := service.GetByName(name)
//
// Create - Create a new resource:
//
//	resource, err := service.Create(data)
//
// Delete - Remove a resource:
//
//	err := service.Delete(id)
//
// Import/Export - Transfer resources between systems:
//
//	exported, err := service.Export(id)
//	imported, err := service.Import(data)
//
// # Pagination
//
// Services handle pagination automatically using QueryParams:
//
//	type QueryParams struct {
//	    Limit int
//	    Skip  int
//	}
//
// The GetAll methods automatically iterate through all pages and return
// the complete result set.
//
// # Thread Safety
//
// Service instances are safe for concurrent use from multiple goroutines.
// Each service maintains no mutable state and delegates to the underlying
// HTTP client which handles connection pooling safely.
package services
