# ipctl

A powerful command-line interface for managing Itential Platform servers. `ipctl` provides intuitive commands for working with 35+ resource types including projects, workflows, automations, adapters, accounts, groups, and more.

## Features

- **Resource Management**: Full CRUD operations (create, read, update, delete) for all Itential Platform resources
- **Import/Export**: Git-based asset import/export with SSH authentication support
- **Multi-Instance Support**: Profile-based configuration for managing multiple Itential Platform servers
- **Flexible Authentication**: OAuth2 client credentials flow or basic authentication
- **Multiple Output Formats**: Human-readable, JSON, YAML, or custom templates
- **Asset Orchestration**: Copy automations, workflows, and other assets between environments
- **Local Development**: Built-in MongoDB-backed AAA server for local testing

## Installation

Download the latest release for your platform from the [releases page](https://github.com/itential/ipctl/releases):

```bash
# Example for Linux x64
curl -LO https://github.com/itential/ipctl/releases/latest/download/ipctl-linux-amd64.tar.gz
tar -xzf ipctl-linux-amd64.tar.gz
sudo mv ipctl /usr/local/bin/
chmod +x /usr/local/bin/ipctl

# Verify installation
ipctl --version
```

## Quick Start

### 1. Configure a Profile

Create a configuration file at `~/.platform.d/config`:

```ini
[profile production]
host = platform.example.com
port = 443
scheme = https
username = admin
password = your-password

[profile staging]
host = staging.example.com
port = 443
scheme = https
client_id = your-client-id
client_secret = your-client-secret
```

### 2. List Resources

```bash
# List all projects on production
ipctl get projects --profile production

# List workflows with JSON output
ipctl get workflows --profile staging --output json

# Get specific automation by name
ipctl get automation "Deploy Network Config"
```

### 3. Export Assets

```bash
# Export a project to local directory
ipctl export project MyProject --destination ./exports/

# Export to Git repository
ipctl export project MyProject \
  --repository git@github.com:org/repo.git \
  --branch main

# Export automation with dependencies
ipctl export automation "My Automation" \
  --include-dependencies \
  --destination ./my-automation/
```

### 4. Import Assets

```bash
# Import from local directory
ipctl import project ./my-project/ --profile production

# Import from Git repository
ipctl import project \
  --repository https://github.com/org/repo.git \
  --branch develop \
  --profile staging
```

### 5. Copy Between Environments

```bash
# Copy automation from staging to production
ipctl copy automation "Deploy Config" \
  --from staging \
  --to production
```

## Configuration

### Configuration File

The configuration file is located at `~/.platform.d/config` and uses INI format:

```ini
[profile default]
host = localhost
port = 3000
scheme = http
username = admin@itential.com
password = admin

[profile production]
host = prod.example.com
port = 443
scheme = https
client_id = your-client-id
client_secret = your-client-secret
verify_ssl = true
```

### Configuration Options

| Option | Description | Required | Default |
|--------|-------------|----------|---------|
| `host` | Platform server hostname | Yes | - |
| `port` | Server port | No | 443 |
| `scheme` | Protocol (http/https) | No | https |
| `username` | Basic auth username | No* | - |
| `password` | Basic auth password | No* | - |
| `client_id` | OAuth2 client ID | No* | - |
| `client_secret` | OAuth2 client secret | No* | - |
| `verify_ssl` | Verify SSL certificates | No | true |

*Either `username`/`password` or `client_id`/`client_secret` required for authentication.

### Environment Variables

Override configuration values using environment variables with the `IPCTL_` prefix:

```bash
export IPCTL_HOST=platform.example.com
export IPCTL_CLIENT_ID=your-client-id
export IPCTL_CLIENT_SECRET=your-client-secret

ipctl get projects
```

### Configuration Precedence

Configuration values are resolved in the following order (highest to lowest):

1. Command-line flags
2. Environment variables (`IPCTL_*`)
3. Configuration file (`~/.platform.d/config`)
4. Default values

## Usage Examples

### Working with Projects

```bash
# List all projects
ipctl get projects

# Get project details
ipctl describe project "My Project"

# Create a new project
ipctl create project "New Project" --description "Project description"

# Delete a project
ipctl delete project "Old Project"

# Export project with Git
ipctl export project "My Project" \
  --repository git@github.com:org/repo.git \
  --branch main \
  --commit-message "Export project"
```

### Working with Automations

```bash
# List automations
ipctl get automations

# Get automation details with JSON output
ipctl describe automation "Deploy Config" --output json

# Import automation from Git
ipctl import automation \
  --repository https://github.com/org/automations.git \
  --reference v1.0.0

# Copy automation between environments
ipctl copy automation "Deploy Config" --from dev --to staging
```

### Working with Workflows

```bash
# List workflows
ipctl get workflows

# Describe workflow with custom template
ipctl describe workflow "My Workflow" --template custom.tmpl

# Export workflow to directory
ipctl export workflow "My Workflow" --destination ./workflows/

# Import workflow
ipctl import workflow ./workflows/my-workflow.json
```

### Working with Adapters

```bash
# List all adapters
ipctl get adapters

# Get adapter details
ipctl describe adapter "ServiceNow"

# Start adapter
ipctl start adapter "ServiceNow"

# Stop adapter
ipctl stop adapter "ServiceNow"

# Restart adapter
ipctl restart adapter "ServiceNow"
```

### Working with Accounts and Groups

```bash
# List accounts
ipctl get accounts

# Create account
ipctl create account "user@example.com" \
  --password "password" \
  --first-name "John" \
  --last-name "Doe"

# List groups
ipctl get groups

# Create group
ipctl create group "Operators" --description "Network operators"

# Add user to group
ipctl update group "Operators" --add-member "user@example.com"
```

### Output Formats

```bash
# Human-readable output (default)
ipctl get projects

# JSON output
ipctl get projects --output json

# YAML output
ipctl get projects --output yaml

# Custom template
ipctl describe project "My Project" --template my-template.tmpl

# No color output
ipctl get projects --no-color
```

## Resource Types

`ipctl` supports management of the following resource types:

### Assets
- Projects
- Workflows
- Automations
- Templates (JSON, Command, Jinja, Analytic)
- Transformations
- Forms
- Operations Manager (OM) instances

### Platform Resources
- Accounts
- Groups
- Roles
- Adapters
- Integrations
- Application instances
- Tags

### Configuration
- Profiles (user profiles, not connection profiles)
- Variables

### Advanced
- Datasets (when feature flag enabled)
- Models
- Device groups

## Development

### Prerequisites

- Go 1.24 or later
- Make
- Git

### Building from Source

```bash
# Clone repository
git clone https://github.com/itential/ipctl.git
cd ipctl

# Install dependencies
make install

# Build
make build

# Binary will be in ./bin/ipctl
./bin/ipctl --version
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific test
go test ./pkg/services/...

# Debug tests
scripts/test.sh debugtest
```

### Project Structure

```
ipctl/
├── cmd/ipctl/           # Application entry point
├── internal/            # Private application code
│   ├── cli/            # CLI setup and command tree
│   ├── handlers/       # Command handlers
│   ├── runners/        # Operation execution
│   ├── flags/          # CLI flag definitions
│   ├── config/         # Configuration management
│   └── profile/        # Profile management
├── pkg/                # Public libraries
│   ├── client/         # HTTP client
│   ├── services/       # API operations
│   ├── resources/      # Business logic
│   └── validators/     # Input validation
├── docs/               # Documentation
└── scripts/            # Build and test scripts
```

For detailed architecture documentation, see [CLAUDE.md](CLAUDE.md).

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Documentation

- [Configuration Reference](docs/configuration-reference.md) - Complete configuration options
- [CLAUDE.md](CLAUDE.md) - Comprehensive architecture and development guide
- [Command Reference](docs/commands.md) - Detailed command documentation (if available)

## Common Issues

### Authentication Failures

```bash
# Verify profile configuration
cat ~/.platform.d/config

# Test connection with verbose output
ipctl get projects --profile myprofile --verbose

# Check environment variables
env | grep IPCTL
```

### SSL Certificate Errors

```bash
# Disable SSL verification (not recommended for production)
ipctl get projects --profile myprofile --verify-ssl=false

# Or set in profile
[profile myprofile]
verify_ssl = false
```

### Profile Not Found

```bash
# List available profiles
grep "^\[profile" ~/.platform.d/config

# Use specific profile
ipctl get projects --profile production

# Set default profile
[profile default]
host = your-default-host.com
```

## License

Copyright 2024 Itential Inc. All Rights Reserved.

Unauthorized copying of this software, via any medium is strictly prohibited. Proprietary and confidential.

## Support

For issues and questions:
- GitHub Issues: https://github.com/itential/ipctl/issues
- Documentation: https://docs.itential.com
