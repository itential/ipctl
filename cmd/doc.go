// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package cmd provides the command-line interface for ipctl.
//
// This package defines the root command structure and coordinates between
// Cobra commands, handlers, and the underlying business logic.
//
// # Command Structure
//
// Commands are organized into logical groups:
//   - Asset Commands: manage projects, workflows, automations, templates, etc.
//   - Platform Commands: interact with the Itential Platform server (start, stop, restart, inspect, API operations)
//   - Dataset Commands: batch operations on assets (load, dump)
//   - Plugin Commands: extended functionality (local-aaa, client)
//
// # Architecture
//
// The command flow follows this pattern:
//
//	User → Cobra CLI → Handler → Runner → Service → Client → API
//
// Each layer has specific responsibilities:
//   - CLI: Parse arguments and flags using Cobra framework
//   - Handler: Coordinate command creation, manage handler registry, format output
//   - Runner: Execute business logic, orchestrate service calls
//   - Service: Implement API interaction logic
//   - Client: Handle HTTP communication with the Itential Platform
//
// # Configuration
//
// The CLI uses profile-based configuration stored in ~/.platform.d/config:
//   - Multiple profiles can be defined for different Itential Platform instances
//   - Configuration precedence: CLI flags → environment variables → config file
//   - Runtime context (client, config, descriptors, verbose flag) is shared across handlers
//
// # Command Descriptors
//
// Commands are defined using YAML descriptors in cmd/descriptors/ that specify:
//   - Command usage, description, and examples
//   - Argument requirements and validation rules
//   - Group assignments for organizing commands in help output
//   - Visibility (hidden) and availability (disabled) flags
//
// # Adding New Commands
//
// To add a new command:
//  1. Define a descriptor YAML file in cmd/descriptors/ or internal/handlers/descriptors/
//  2. Create a runner in internal/runners/ implementing the necessary interfaces (Reader, Writer, etc.)
//  3. Create a handler in internal/handlers/ that wraps the runner with an AssetHandler
//  4. Register the handler in handlers.NewHandler()
//  5. The handler registry automatically creates commands based on implemented interfaces
//
// # Example
//
//	// Create and execute the CLI
//	exitCode := cmd.Execute()
//	os.Exit(exitCode)
package cmd
