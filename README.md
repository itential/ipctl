# ipctl

The `ipctl` project provides a command line interface for working with Itential Platform servers.

## Getting Started

Getting starting using `ipctl` is straightforward and easy.  This section will
guide you through getting up and running.

### Install `ipctl`

To install `ipctl` navigate to the [releases page](https://github.com/itential/ipctl/releases) and download the tarball that matches your machine OS and architecture.  Then  untar the compressed file put the `ipctl` executable in your path.

### Configuring `ipctl`

The default configuration file for `ipctl` is `~/.platform.d/config`. This file
is read by `ipctl` and used to configure various options.  The configuration
file is divided into multiple sections, each section providing one or more
key=vlaue configuration entries.

The configuration file can be configured using any standard text editor and
should be secured such that it is only readable by the currently logged in
user.

The `ipctl` application can also be configured using environment variables.  If
configured, environment variables take precedence over values set in the
configuration file.

The configuration file uses profiles to define one or more target instances of Itential Platform.   For example, the following will create a new profile called `devel` in the configuration file.

```
[profile devel]
host = devel.itential.com
```

Once created, use the profile name when running commands.

```
ipctl export project test --profile devel
```

More than one profile can be created, each with a different set of parameters.  The following options can be configured in a profile.

- host - Configures the hostname or IP address of the Itential Platform server
- port - Configures the port to use when connecting to the server  The default value is either `80` when `use_tls=false` or `443` when `use_tls=true`
- use_tls - Boolean value that enables or disables the use of TLS when connecting to the server (default=`true`)

- username - Configures the username to use when authenticating to Itential Platform (default=`admin@pronghorn`)
- password - Conifugres the password to use when authenticating to Itential Platform (default=`admin`)

- client_id - Conifigures the client id value when using an Itential Cloud service account
- client_secret - Conifgures the client secret when using an Itential Cloud service account

Note: When both `client_id` and `client_secret` are configured, the values for `username` and `password` are ignored by the CLI.

If you do not want to pass the `--profile` option for every command, you can set the name of the default profile.  To set the name of the default profile add the following to your configuration flie.

```
[application]
default_profile = <name_of_default_profile_to_use>
```

Note: The profile must be configured for `default_profile` to work

Additional logging can be enabled in the configuration file as well.   To enable `debug` logging, add the following configuration block.

```
[log]
level = debug
```


## Configuring profiles

The `ipctl` configuration file supports configuring one more more profiles.  Profiles define the properties for connecting to an Itential Platform server instance.   The following configuration settings can be used to configure a profile.

Profile define connection properties for connecting to an instnace of Itential
Platform.  It includes the connection parameters, authentication parameters and session parameters.   Any value configured in a profile will override the same value set in `defaults`.

Any profile setting can also be set using an environment variable.  The override a profile setting with an environment variable, use the form of `IPCTL_PROFILE_<PROFILE>_<KEY>=<VALUE>`.

For instance, assume passing the value of password using an environment variable as opposed to putting it into the configuration file.  In order to set the value of password for a profile call `prod`, the environment variable would be `IPCTL_PROFILE_PROD_PASSWORD=itsasecret`

### `host`

The `host` setting defines the hostname or IP address of the target Itential Platform server.  The default value is `localhost`

### `port`

The `port` setting configurees the port to use when connecting to the serer.  The  default value for `port` is determined based on the value for `use_tls`.  When `use_tls` is set to `true` (default), the default port is `443`.  The `use_tls` is set to `false`, the default port is `80`.

### `use_tls`

The `use_tls` setting is a boolean setting that enables or disables the use of TLS connections to the Itential Platform server.   When the `use_tls` value is set to `true` the client will attept to use TLS when connecting to the server.  When `use_tls` is set to `false`, the client will not use TLS when connecting to the server.

The default value for `use_tls` is `true`

### `verify`

The `verify`  setting enables or disables the certificate verification.  When connecting to an Itential Platform server using TLS, the client will attempt to verify the server certificate.  In some cases (such as with self signed certificates), the client should not attempt to verify the certificate.  Setting this value to `false` will disable certificate verification.

The default value for `verify` is `true`

### `usernam`

The `username` setting configures the name of the account to use when authenticating to an Itential Platform server.

The default value for `username` is `admin@pronghorn`

### `password`

The `password` setting configures the password to use when authenticating the account.  It is used in conjunction with the `username` setting.

The default `password` is `admin`

### `client_id`

### `client_secret`

### `timeout`

### `output`

the `output` setting defines the default output to return when running CLI
commands.  This setting accepts one of `human` or `json`.  When set to `human`,
the output is returned in human readable format.  When set to `json`, the
native JSON output is returned.

The default value for `output` is `human`

### `verbose`

### `pager`

The `pager` setting enables or disables the the pager output.  When this value
is enabled, the output from commands are piped through `less` for pagination.
When the value is set to `false`, the command is returned directly to `stdout`.

The default value for `pager` is `true`



