// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"embed"

	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/pkg/client"
)

//go:embed templates/*
var templates embed.FS

// BaseRunner provides common fields and functionality for all runners.
// All specific runners should embed BaseRunner to inherit shared configuration
// and client access.
//
// By using config.Provider interface instead of *config.Config, runners are
// decoupled from the concrete configuration type, making them easier to test
// and more flexible.
//
// Example usage:
//
//	type ProjectRunner struct {
//	    BaseRunner
//	    service      *services.ProjectService
//	    accounts     *services.AccountService
//	}
//
//	func NewProjectRunner(client client.Client, cfg config.Provider) *ProjectRunner {
//	    return &ProjectRunner{
//	        BaseRunner: NewBaseRunner(client, cfg),
//	        service:    services.NewProjectService(client),
//	        accounts:   services.NewAccountService(client),
//	    }
//	}
type BaseRunner struct {
	// config provides access to application configuration via interfaces.
	// Use specific interfaces (ProfileProvider, GitProvider, etc.) when possible
	// for better encapsulation.
	config config.Provider

	// client provides HTTP access to the Itential Platform APIs.
	// Use this to create service instances as needed.
	client client.Client
}

// NewBaseRunner creates a new BaseRunner with the provided client and configuration.
// This should be called from specific runner constructors.
//
// The configuration is accepted as a Provider interface rather than *config.Config,
// which enables dependency injection and makes testing easier.
func NewBaseRunner(client client.Client, cfg config.Provider) BaseRunner {
	return BaseRunner{
		config: cfg,
		client: client,
	}
}

// Config returns the configuration provider.
// This returns the Provider interface, allowing access to all configuration aspects.
// For better encapsulation, consider accessing specific methods directly rather than
// storing the entire config.
func (r *BaseRunner) Config() config.Provider {
	return r.config
}

// Client returns the HTTP client for creating services.
func (r *BaseRunner) Client() client.Client {
	return r.client
}
