// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"embed"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
)

//go:embed templates/*
var templates embed.FS

// BaseRunner provides common fields and functionality for all runners.
// All specific runners should embed BaseRunner to inherit shared configuration
// and client access.
//
// Example usage:
//
//	type ProjectRunner struct {
//	    BaseRunner
//	    service      *services.ProjectService
//	    accounts     *services.AccountService
//	}
//
//	func NewProjectRunner(client client.Client, cfg *config.Config) *ProjectRunner {
//	    return &ProjectRunner{
//	        BaseRunner: NewBaseRunner(client, cfg),
//	        service:    services.NewProjectService(client),
//	        accounts:   services.NewAccountService(client),
//	    }
//	}
type BaseRunner struct {
	// config contains the global configuration including profiles, logging,
	// and terminal settings.
	config *config.Config

	// client provides HTTP access to the Itential Platform APIs.
	// Use this to create service instances as needed.
	client client.Client
}

// NewBaseRunner creates a new BaseRunner with the provided client and configuration.
// This should be called from specific runner constructors.
func NewBaseRunner(client client.Client, cfg *config.Config) BaseRunner {
	return BaseRunner{
		config: cfg,
		client: client,
	}
}

// Config returns the global configuration.
func (r *BaseRunner) Config() *config.Config {
	return r.config
}

// Client returns the HTTP client for creating services.
func (r *BaseRunner) Client() client.Client {
	return r.client
}
