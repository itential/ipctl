// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package config provides configuration management for the ipctl CLI application.
//
// This package handles loading, parsing, and managing configuration settings
// from multiple sources including config files, environment variables, and
// command-line flags. It supports multiple connection profiles for different
// Itential Platform instances.
//
// # Configuration Sources
//
// Configuration values are resolved with the following precedence (highest to lowest):
//
//  1. Command-line flags (--profile, --config, etc.)
//  2. Environment variables (IPCTL_*)
//  3. Configuration file (~/.platform.d/config)
//  4. Default values
//
// # Configuration File
//
// The default configuration file location is:
//
//	~/.platform.d/config
//
// The file uses INI format with sections for each profile:
//
//	[application]
//	working_dir = ~/.platform.d
//	default_profile = production
//
//	[log]
//	level = INFO
//	file_enabled = false
//
//	[terminal]
//	no_color = false
//	default_output = human
//	pager = false
//
//	[profile default]
//	host = localhost
//	port = 3000
//	use_tls = true
//	verify = false
//	username = admin@pronghorn
//	password = admin
//	timeout = 30
//
//	[profile production]
//	host = platform.example.com
//	port = 443
//	use_tls = true
//	verify = true
//	client_id = your-client-id
//	client_secret = your-client-secret
//
// # Profiles
//
// Profiles define connection settings for different Itential Platform instances.
// Each profile includes:
//
//   - Host: Server hostname or IP address
//   - Port: Server port number
//   - UseTLS: Whether to use HTTPS
//   - Verify: Whether to verify TLS certificates
//   - Authentication: Username/password or OAuth2 credentials
//   - Timeout: Request timeout in seconds
//
// # Creating Configuration
//
// Initialize configuration with defaults:
//
//	cfg := config.NewConfig(nil, nil, "", "", "")
//
// Or with custom settings:
//
//	defaults := map[string]interface{}{
//	    "application.working_dir": "~/.myapp",
//	}
//	envBindings := map[string]string{
//	    "application.working_dir": "MYAPP_WORKING_DIR",
//	}
//	cfg := config.NewConfig(defaults, envBindings, "~/.myapp", "", "config")
//
// # Accessing Profiles
//
// Get the active profile:
//
//	profile, err := cfg.ActiveProfile()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Get a specific profile by name:
//
//	profile, err := cfg.GetProfile("production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Profile Structure
//
// A Profile contains connection settings:
//
//	type Profile struct {
//	    Host         string
//	    Port         int
//	    UseTLS       bool
//	    Verify       bool
//	    Username     string
//	    Password     string
//	    ClientID     string
//	    ClientSecret string
//	    MongoUrl     string
//	    Timeout      int
//	}
//
// # Environment Variables
//
// All configuration values can be overridden with environment variables:
//
// Application settings:
//   - IPCTL_APPLICATION_WORKING_DIR
//   - IPCTL_APPLICATION_DEFAULT_PROFILE
//   - IPCTL_APPLICATION_DEFAULT_REPOSITORY
//
// Feature flags:
//   - IPCTL_FEATURES_DATASETS_ENABLED
//
// Logging:
//   - IPCTL_LOG_LEVEL (DEBUG, INFO, WARN, ERROR)
//   - IPCTL_LOG_FILE_JSON (true/false)
//   - IPCTL_LOG_CONSOLE_JSON (true/false)
//   - IPCTL_LOG_FILE_ENABLED (true/false)
//   - IPCTL_LOG_TIMESTAMP_TIMEZONE (utc/local/timezone)
//
// Terminal:
//   - IPCTL_TERMINAL_NO_COLOR (true/false)
//   - IPCTL_TERMINAL_DEFAULT_OUTPUT (human/json/yaml)
//   - IPCTL_TERMINAL_PAGER (true/false)
//
// Git:
//   - IPCTL_GIT_NAME
//   - IPCTL_GIT_EMAIL
//   - IPCTL_GIT_USER
//
// Profile-specific (replace PROFILE with profile name):
//   - IPCTL_PROFILE_PROFILE_HOST
//   - IPCTL_PROFILE_PROFILE_PORT
//   - IPCTL_PROFILE_PROFILE_USE_TLS
//   - IPCTL_PROFILE_PROFILE_VERIFY
//   - IPCTL_PROFILE_PROFILE_USERNAME
//   - IPCTL_PROFILE_PROFILE_PASSWORD
//   - IPCTL_PROFILE_PROFILE_CLIENT_ID
//   - IPCTL_PROFILE_PROFILE_CLIENT_SECRET
//   - IPCTL_PROFILE_PROFILE_TIMEOUT
//
// # Command-line Flags
//
// Override configuration with flags:
//
//	ipctl --profile production --config /path/to/config get projects
//
// The --profile flag selects which profile to use for the command.
// The --config flag specifies an alternate configuration file.
//
// # Repositories
//
// Git repositories can be configured for import/export operations:
//
//	[repository myrepo]
//	url = git@github.com:user/repo.git
//	reference = main
//	private_key_file = ~/.ssh/id_rsa
//
// Access repositories by name:
//
//	repo, err := cfg.GetRepository("myrepo")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Configuration Dumping
//
// Debug configuration by dumping current values:
//
//	fmt.Println(cfg.DumpConfig())
//
// This outputs JSON representation of all loaded configuration.
//
// # Thread Safety
//
// Config instances are safe for concurrent read access but should not be
// modified concurrently. Initialize configuration once at startup and treat
// it as immutable during runtime.
//
// # Example: Complete Setup
//
//	// Initialize configuration
//	cfg := config.NewConfig(nil, nil, "", "", "")
//
//	// Get active profile
//	profile, err := cfg.ActiveProfile()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Create HTTP client
//	client, err := client.NewHttpClient(profile)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Use client with services
//	projectSvc := services.NewProjectService(client)
//	projects, err := projectSvc.GetAll()
//	if err != nil {
//	    log.Fatal(err)
//	}
package config
