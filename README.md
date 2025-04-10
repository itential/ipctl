# ipctl

The `ipctl` project provides a command line interface for working with Itential Platform servers.

## Getting Started

Getting starting using `ipctl` is straightforward and easy.  This section will
guide you through getting up and running.

### Install `ipctl`

To install `ipctl` navigate to the [releases page](https://github.com/itential/ipctl/releases) and download the tarball that matches your machine OS and architecture.  Then untar the compressed file put the `ipctl` executable in your path.

### Configuring `ipctl`

The default configuration file for `ipctl` is `~/.platform.d/config`. This file
is read by `ipctl` and used to configure various options.  The configuration
file is divided into multiple sections, each section providing one or more
key=value configuration entries.

The configuration file can be configured using any standard text editor and
should be secured such that it is only readable by the currently logged in
user.

The `ipctl` application can also be configured using environment variables.  If
configured, environment variables take precedence over values set in the
configuration file.

The configuration file uses profiles to define one or more target instances of Itential Platform.  For example, the following will create a new profile called `devel` in the configuration file.

```ini
[profile devel]
host = devel.itential.com
```

Once created, use the profile name when running commands.

```bash
ipctl export project test --profile devel
```

For a complete list of all configuration parameters, please see the
configuration reference [here](docs/configuration-reference.md)
