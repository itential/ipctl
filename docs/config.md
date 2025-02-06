# Configuration Reference.

This document provides details about the availalbe configuration options
available to configure `ipctl`.   All configuration options can be specified
either in a configuration file and/or as environment variables.

Values configured using environment variables take precedence over the same
values configured in a configuration file.

By default, `ipctl` will load the configuration file `~/.platform.d/config`.
To use a different configuraiton file name and/or location, use the `--config`
command line option and pass it the path to the configuration file you wish to
use.

The configuration file can also be specified using the `IPCTL_CONFIG_FILE`
environment variable.

The configuration file is organzied into sections with each section providing
one or more key=value settings that can be used to configure the application.

## Avoiding security risks

The `ipctl` configuration file may contain senstive configuration values used
to authenticate to services.   The configuration file should be secured such
that only the current user has access to the configuration directory and file.

The following example commands would security the directory and file when
running `ipctl` from a Linux host.

```bash
chmod 700 ~/.platform.d
chmod 600 ~/.platform.d/config
```

Please consult your OS documentation for specific commands.

## Using profiles

The application configuration flie supports configuring multiple profiles  A
profile defines the connection settings for a given server.  When `iapctl`
attempts to conect to a server, it will look up the profile settings in the
conifguraiton file based on the profile name.

The profile name can be passed to any command using the `--profile <name>`
command line argument.  A default profile can also be set in the configuration
file.  See `default_profile` in the [Applicaiton Settings}(#Application
Settings) section for a descrption.

## Configuration options

The entries below provide the set of available configuration options that can
be used with`ipctl`.

### Application Settings

The `application` settings can be used to configure application level settings
for `ipctl`.  The following values are configurable within `[application]`

#### `default_profile`

Configures the name of the profile to use if not explicitly set using the
command line option `--profile`

The default value for `default_profile` is `null`

#### working_dir

### Log Settings

#### level

#### file_json

#### console_json

#### file_enabled

#### timestamp_timzezone

### Terminal Settings

#### no_color

#### timestamp_timezone

### Profiles

#### `host`

Sets the hostname or IP address of the Itential Platform server to connect to.

The default value for `host` is `localhost`

#### `port`

Sets the `port` value to use when connecting to the Itential Platform server.
This setting should be a numeric value used to connect to the server.  When the
value of `port` is set to `0`, the value is automatically determined based on
the value of `use_tls`.  When `use_tls` is `true` the value of `port` will be
set to `443`.  When `use_tls` is set to `false`, the value of `port` will be
set to `80`.

The default value for `port` is `0`.

#### `use_tls`

Enables or disables the use of TLS for the connecton.  When this value is set
to `true`, the application will attempt to establish a TLS connecton to the
server.   When this value is set to `false` the application will not attempt to
use TLS when connecting to the server.

The default value for `use_tls` is `true`


#### `verify`

Enables or disables certificate validation for TLS based connections.  When
this value is set to `true`, certificates received from the server are
validated. When this value is set to `false`, certifcates are assumed valid.
This feature is useful when using self-signed certificates for TLS connections.

The default value for `verify` is `true`

#### `username`

Cnfigures the name of the user to use when authenticating to the Itential
Platform server.  This should be the username used to login to the server and
determines the level of authorzation.

The default value for `username` is `admin@pronghorn`

#### `password`

Configures the password to use when authenticating the connection to Itential
Platform server.

The default value for `password` is `admin`

#### client_id

#### client_secret

#### timeout

#### output

#### pager


### Repository settings

#### url

#### private_key

#### private_key_file

#### reference

#### name

#### email

### Mongo Settings

#### url








