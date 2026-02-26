# 💻 ipctl

A CLI for managing Itential Platform servers — get, create, import, export, and copy 35+ resource types from the command line.

![ipctl version](assets/gifs/ipctl-version.gif)

## Quick Start

### Install

Download the latest release for your platform from the [releases page](https://github.com/itential/ipctl/releases):

```bash
# Linux x64
curl -LO https://github.com/itential/ipctl/releases/latest/download/ipctl-linux-x86_64.tar.gz
tar -xzf ipctl-linux-x86_64.tar.gz
sudo mv ipctl /usr/local/bin/
```

<details>
<summary>Build from source</summary>

```bash
git clone https://github.com/itential/ipctl.git && cd ipctl
make build
./bin/ipctl --version
```

</details>

### Configure

Create `~/.platform.d/config.toml` with a server profile:

```toml
["profile default"]
host = "platform.example.com"
port = 443
use_tls = true
username = "admin"
password = "your-password"
```

### Use

![ipctl workflow](assets/gifs/ipctl-workflow.gif)

## Features

| Feature | Description |
|---------|-------------|
| **Resource CRUD** | Get, create, update, and delete 35+ resource types |
| **Import / Export** | Move assets via local directories or Git repositories with SSH auth |
| **Multi-Instance** | Named profiles for managing multiple Platform servers |
| **Authentication** | OAuth2 client credentials or basic auth with TLS |
| **Output Formats** | Human-readable tables, JSON, YAML, or custom Go templates |
| **Cross-Environment** | Copy automations, workflows, and assets between servers |

## Configuration

`ipctl` loads configuration from `~/.platform.d/config` by default. Supports INI, YAML, TOML, and JSON formats (auto-detected by file extension).

### Profile Options

| Option | Description | Default |
|--------|-------------|---------|
| `host` | Platform server hostname | `localhost` |
| `port` | Server port (0 = auto from `use_tls`) | `0` |
| `use_tls` | Enable TLS connection | `true` |
| `verify` | Verify TLS certificates | `true` |
| `username` | Basic auth username | - |
| `password` | Basic auth password | - |
| `client_id` | OAuth2 client ID | - |
| `client_secret` | OAuth2 client secret | - |
| `timeout` | Request timeout in seconds (0 = disabled) | `0` |

Authentication requires either `username`/`password` or `client_id`/`client_secret`.

### Environment Variables

Override any profile value with `IPCTL_PROFILE_<NAME>_<KEY>`:

```bash
export IPCTL_PROFILE_PROD_PASSWORD=secret
ipctl get projects --profile prod
```

### Precedence

CLI flags > environment variables > config file > defaults

See the [Configuration Reference](docs/configuration-reference.md) for complete details including multi-format examples.

## Supported Resources

| Category | Resources |
|----------|-----------|
| **Automation Studio** | projects, workflows, automations, templates, transformations, jsonforms |
| **Admin** | accounts, groups, roles, adapters, integrations, prebuilts, tags |
| **Configuration Manager** | devices, device-groups, configuration-parsers, gctrees |
| **Lifecycle Manager** | models |

See the [Command Quick Reference](docs/commands-quick-reference.md) for the full matrix of supported operations per resource.

## Usage Examples

<details>
<summary>Working with projects</summary>

```bash
ipctl get projects
ipctl describe project "My Project"
ipctl create project "New Project" --description "Project description"
ipctl delete project "Old Project"
ipctl export project "My Project" \
  --repository git@github.com:org/repo.git \
  --branch main
```

</details>

<details>
<summary>Working with automations</summary>

```bash
ipctl get automations
ipctl describe automation "Deploy Config" --output json
ipctl import automation \
  --repository https://github.com/org/automations.git \
  --reference v1.0.0
ipctl copy automation "Deploy Config" --from dev --to staging
```

</details>

<details>
<summary>Working with adapters</summary>

```bash
ipctl get adapters
ipctl describe adapter "ServiceNow"
ipctl start adapter "ServiceNow"
ipctl stop adapter "ServiceNow"
ipctl restart adapter "ServiceNow"
```

</details>

<details>
<summary>Output formats</summary>

![ipctl output formats](assets/gifs/ipctl-output-formats.gif)

</details>

## Documentation

- [Configuration Reference](docs/configuration-reference.md) — profile options, formats, environment variables
- [Command Quick Reference](docs/commands-quick-reference.md) — operations matrix per resource
- [API Command Reference](docs/api-command-reference.md) — detailed API command docs
- [Working with Repositories](docs/working-with-repositories.md) — Git-based import/export
- [Logging Reference](docs/logging-reference.md) — log levels, JSON output, sensitive data redaction
- [Running from Source](docs/running-from-source.md) — development setup

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, branch conventions, and the pull request process.

All contributors must sign the [Contributor License Agreement](CLA.md) before contributions can be merged.

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE).


## Support

- [GitHub Issues](https://github.com/itential/ipctl/issues)
- [Itential Documentation](https://docs.itential.com)
