# Configuration Reference

This document provides details about the available configuration options
available to configure `ipctl`.  All configuration options can be specified
either in a configuration file and/or as environment variables.

Values configured using environment variables take precedence over the same
values configured in a configuration file.

By default, `ipctl` will load the configuration file `~/.platform.d/config`.
To use a different configuration file name and/or location, use the `--config`
command line option and pass it the path to the configuration file you wish to
use.

The configuration file can also be specified using the `IPCTL_CONFIG_FILE`
environment variable.

The configuration file is organized into sections with each section providing
one or more key=value settings that can be used to configure the application.

## Avoiding security risks

The `ipctl` configuration file may contain sensitive configuration values used
to authenticate to services.  The configuration file should be secured such
that only the current user has access to the configuration directory and file.

The following example commands would secure the directory and file when
running `ipctl` from a Linux host.

```bash
chmod 700 ~/.platform.d
chmod 600 ~/.platform.d/config
```

Please consult your OS documentation for specific commands.

## Using profiles

The application configuration file supports configuring multiple profiles.  A
profile defines the connection settings for a given server.  When `ipctl`
attempts to connect to a server, it will look up the profile settings in the
configuration file based on the profile name.

The profile name can be passed to any command using the `--profile <name>`
command line argument.  A default profile can also be set in the configuration
file.  See `default_profile` in the [Application Settings}(#Application
Settings) section for a description.

## Configuration options

The entries below provide the set of available configuration options that can
be used with `ipctl`.

### Application Settings

The `application` settings can be used to configure application level settings
for `ipctl`.  The following values are configurable within `[application]`.
All values in this section can be overridden with environment variables.

For instance, to override the `default_profile` value, the environment variable
would be `IPCTL_APPLICATION_DEFAULT_PROFILE`.

#### `default_profile`

Configures the name of the profile to use if not explicitly set using the
command line option `--profile`.

The default value for `default_profile` is `null`.

#### working_dir

Configures the working directory for the application.  The working directory
is the default directory where the application will look for the configuration
file.

The default value for `working_dir` is `~/.platform.d`.

### Log Settings

The log settings section exposes configuration settings that can be used to
change how the application logging is performed.  All configuration values are
defined under the `[log]` section.   Any values configured in this section can
be overridden using environment variables in the form of `IPCTL_LOG_<name>`
where `<name>` is the key name.

#### level

Configures the level of detail sent to the logging facility.  This
configuration value accepts one of the following values: `info`, `debug`,
`trace`.

The default value for `level` is `info`.

#### file_json

#### console_json

#### file_enabled

#### timestamp_timzezone

### Terminal Settings

The terminal settings section provides configuration values for managing the
terminal environment.  All configuration settings are maintained under the
`[terminal]` section and can be overridden using environment variables
prefaced with `IPCTL_TERMINAL_<NAME>`.

#### no_color

Enables or disables the use of color in `human` output from the application.
When this value is set to `true`, no output is colorized and when this value
is `false`, the output may use color in the output.

This configuration value only applies to `human` output.

The default value for `no_color` is `false`.

#### default_output

Sets the default output format for commands.  Currently, the application
supports three output formats `human`,`json` and `yaml`.   Use this
configuration to define the default output format for all commands.

This setting can be override for any command using `--output <format>`.

The default value for `default_output` is `human`.

#### pager

Enables or disables the pager feature in the application.  The pager feature
will pass the returned output through `less`, which must be available on your
system, to paginate the output.

The default value for `pager` is `true`.

#### timestamp_timezone

Configures the timezone to use when converting log timestamp messages from the
application.  This setting can be used to automatically translate the log
messages to any desired timezone.

The default value for `timestamp_timezone` is `utc`.

### Git settings

When working with git repositories, there are some global settings that should
be configured.  This section provides a configuration settings for setting
global git options.

#### `name`

Configures the name to use when making commits using git.  The name configured
here will be displayed in the commit message.

The default value for `name` is `null`.

#### `email`

Configures the email address to use when making git commits.  The email address
will be used in commit messages.

The default value for `email` is `null`.

#### `user`

Sets the default user to use when connecting to a git repository over SSH.
Most git server implementations require this value to be `git` which is the
default; however it can be changed here if needed.

The default value for `user` is `git`.

### Profiles

Profiles provide server specific configuration settings for working with
Itential Platform servers.  The configuration file format supports creating one
or more named profiles.

To create a named profile, start the section with `[profile <name>]`.  For
instance, to create a new provide called `staging`, the section would be
`[profile staging]`.  Once created, the profile can be invoked by name to the
`--profile` command line argument.

There is one special profile called `default`.  The `default` profile can be
used to configure profile default values that will be applied to every profile
unless specifically overridden.

All profile settings can also be overridden using environment variables.  For
instance, assume we want to pass the password in using an environment variable
instead of storing the password in the configuration file.  This could be
accomplished by setting `IPCTL_PROFILE_<NAME>_PASSWORD` to the desired value.

For instance, to set the password for a profile called `prod`, the environment
variable would be `IPCTL_PROFILE_PROD_PASSWORD`.  The value would override any
value in the configuration file.

#### `host`

Sets the hostname or IP address of the Itential Platform server to connect to.

The default value for `host` is `localhost`.

#### `port`

Sets the `port` value to use when connecting to the Itential Platform server.
This setting should be a numeric value used to connect to the server.  When the
value of `port` is set to `0`, the value is automatically determined based on
the value of `use_tls`.  When `use_tls` is `true` the value of `port` will be
set to `443`.  When `use_tls` is set to `false`, the value of `port` will be
set to `80`.

The default value for `port` is `0`.

#### `use_tls`

Enables or disables the use of TLS for the connection.  When this value is set
to `true`, the application will attempt to establish a TLS connection to the
server.  When this value is set to `false` the application will not attempt to
use TLS when connecting to the server.

The default value for `use_tls` is `true`.

#### `verify`

Enables or disables certificate validation for TLS based connections.  When
this value is set to `true`, certificates received from the server are
validated. When this value is set to `false`, certificates are assumed valid.
This feature is useful when using self-signed certificates for TLS connections.

The default value for `verify` is `true`.

#### `username`

Configures the name of the user to use when authenticating to the Itential
Platform server.  This should be the username used to login to the server and
determines the level of authorization.

The default value for `username` is `admin@pronghorn`.

#### `password`

Configures the password to use when authenticating the connection to Itential
Platform server.

The default value for `password` is `admin`.

#### client_id

Sets the client identifier for working with Itential Platform running in the
cloud at itential.io.  The client id can be obtained when creating a service
account in Itential Cloud.

The default value for `client_id` is `null`.

#### client_secret

Configures the client secret for working with Itential Platform running in
Itential Cloud.  The client secret can be obtained when creating a service
account in Itential Cloud.

The default value for `client_secret` is `null`.

#### timeout

Configure the timeout value in seconds for request messages sent to the
server.

The default value for `timeout` is `5`.

#### mongo_url

Configures the URL to use to make calls directly to the Mongo database.  This
is primarily used by plugins.

The default value for `mongo_url` is `null`.

### Repository settings

Repositories allow for the configuration of named repository configurations for
working with the `push` and `pull` commands.  The configuration file supports
creating one or more repository configuration entries in the configuration
file.

Repository configurations are sections prefixed with `repository`.  For
instance, to create a repository named `assets`, the section would be
`[repository assets]`.  Once configured, the repository can be referenced by
name in various commands.

Any configuration value included in the `repository` section can be overridden
with environment variables prefixed with `IPCTL_REPOSITORY_<NAME>_<KEY>`.  For
instance, to override the `reference` configuration value for the repository
named `assets`, the environment variable would be
`IPCTL_REPOSITORY_ASSETS_REFERENCE`.

#### url

Provides the full URL to the git repository to use for this named repository
object.  This configuration setting accepts any valid git URL format including:

- SSH Transport URLs
- Git Transport URLs
- HTTP/S Transport URLs
- Local Transport URLs

The default value for `url` is `null`.

#### private_key

Configures the private key to use when authenticating to the git repository
over SSH.  This value contains the actual private key.

Note: This configuration setting is mutually exclusive with `private_key_file`.

The default value for `private_key` is `null`.

#### private_key_file

Configures the path to the file containing the private key to use when
authenticating to the repository over SSH.

Note: This configuration setting is mutually exclusive with `private_key`.

The default value for `private_key_file` is `null`.

#### reference

Configures the git reference to use with this named repository.  The git
reference can be a git branch, tag or specific SHA.  If this value is not
configured, the default branch configured on the server will be used.

The default value for `reference` is `null`.

#### name

Configures the name to be included in commit messages when using the
`push` command.

The default value for `name` is `null`.

#### email

Configures the email address to include in commit messages created when using
the `push` command.

The default value for `email` is `null`.
