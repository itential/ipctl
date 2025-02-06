# Working with repositories

One of the key features of `ipctl` is the opportunity to work natively with
`git` repositories.   This integration allows you to directly push and pull
Itential Platform assets to and from git repositories.

The `ipctl` application provides two commands for working with git
repositories, `push` and `pull`.   The `ipctl push` command is used to export
an asset from an Itential Platform server to a git repository.  The `ipctl pull`
command is used to pull an asset from a git repository and import it to the
Itential Platform server.

Before using `push` and `pull`, the configuration file must be updated to
include one more more repository entries.   By default, the configuration file
is located at `~/.platform.d/config`.

## Configuring named repositories

To add a repository to the configuraiton, edit the configuration file with a
text editor and add a `[repository <name>]` section.  For instance, let's
assume we wish to create a new repository configuration called `assets`.

```
[repository assets]
url = https://github.com/itential/assets
```

If working with private or protected repositories, a SSH URL and private key
can be specified.

```
[repository assets]
url = git@github.com:itential/assets.git
private_key_file = ~/.ssh/id_ecdsa
```

Both examples above will configure a new named repository called `assets` that
can be used to push and pull documents to and from Itential Platform servers
and git.

## Working with `push`

The `push` command is used to export (push) an asset from Itential Platform
and push it to a git repository.   This command will attempt to export a named
asset of a specified type and commit it to a named repository.

The general form of the command is `ipctl push <asset> <name> <repo> [options]`
where `asset` is the type of asset (e.g. workflow, transformation, etc), `name`
is the name of the asset as it appears in the Itential Platform UI and `repo`
is the name of the repository as specified in the configuration file.

Below is a very basic example of this command.

```
ipctl push project "Firewall upgrades" assets
```

The above example command will attempt to export the `Firewall upgrades`
project from Itential Platform and add (or update) it to the `assets` named
repository.

The project fill `"Fireawll upgrades.project.json"` will be placed at the root
of the repository.   If the project file does not exist in the repository, it
will be added as a new file.   If the project file does exist in the repository
and has been changed, it will be updated.  If the project file does exist and
is has not been updated, no change will be pushed to the repository.

### Adding a custom commit message

By default, `ipctl` will add a default commit message for every push into a
repository.   A customized commit message can be included with each push using
the `--message` command line option.

For instance, assume we want to add a custom message when committing a
particular workflow, the command would be:

```
ipctl push workflow "Example workflow" assets --message "This is a custom commit message"
```

When the commit is made to the git repository, the commit message will be as
shown above instead of the default commit message.

### Using a different git reference

The git reference can be used to specify the specific branch to use when
committing the asset to the git repository.  The reference can be configured in
the configuration file.  If the reference is not specified, the default branch
as specified by the server is used.

The reference can also be set when performing the commit using the command line
option.  For instance, assume when pushing a new commit, we want to push the
commit to the `staging` branch instead of the default branch.

```
ipctl push workflow "Example workflow" assets --reference staging
```

Note: The referenced branch must already exist on the git server

### Setting the destination path

As previously noted, when pushing an asset into a repository, by default, the
asset is placed in the root of the repository.  In some cases it is deseriable
to build a folder structure within the repository.  To place an asset into a
folder, using the `--path` command line option.

For example, assume we want to push a workflow called `Example workflow` into a
folder path in the repo that is `/staging/workflows`.  The command would be

```
ipctl push workflow "Example workflow" assets --path /staging/workflows
```

if the filepath specified by `--path` does not exist within the directory
structure of the repository, it will be created.

## Working with `pull`

The `pull` command can be used to import (pull) an asset from a name git
repository and import it into the Itential Platform server.  This command will
import the asset to Itential Platform drectly from the git repository.

The general from of the command is `ipctl pull <asset> <filename> <repo>
[options]` where `asset` is the type of asset (e.g. workflow, transformation,
etc), `filename` is the filename of the asset in the repository and `repo` is
the name of the repository as specified in the configuration file.

Below is a simple example of pulling an asset from a repository.

```
ipctl pull workflow "Test workflow.workflow.json" assets
```

The above command will attempt to pull the file `"Test workflow.workflow.json"`
from the named repository called `assets` and import it into the Itential
Platform server.

If a workflow by that name already exists in Itential Platform or the specified
filename cannot be found, this command will return an error.

### Replacing an existing asset

By default, when the `pull` command runs, if the asset name already exists in
the destination Itential Platform server, the command will return an error
saying the asset already exists.   Sometimes it is desirable to replace the
asset on the server with the one from the named git repository.

In order to replace the asset on the server, use the `--replace` command line
option.

```
ipctl pull workflow "Test workflow.workflow.json" assets --replace
```

The example command above will first check if `"Test workflow"` exists on the
desitnation server.  If it does, it will delete the asset on the server and
then proceed to import the asset from the repository.

If the workflow does not exist on the destination server, the workflow simply
be imported.

### Working with references

When pulling an asset from a name git repository, sometimes it is beneficial to
pull from a specific git reference (e.g. tag, branch, sha, etc) that is
different from either the default one or the one specified in the configuration
file.

To pull an asset from a specific reference, use the `--reference` comamnd line
option.  For instance, assume we want to pull the workflow named `"Test
workflow.workflow.json"` from the `staging` branch.

```
ipctl pull workflow "Test workflow.workflow.json" assets --reference staging
```

The command above will attempt to pull the workflow from the branch named
`staging` instead of the configured or default branch.

### Specifying the path

By default, the `pull` command will attempt to find the filename in the root of
the repository.  In some cass, the file could be nested in a directory
structure.  In order to inform the `pull` command as to the location of the
filename, use the `--path` command line option.

```
ipctl pull workflow "Test workflow.workflow.json" assets --path /workflows
```

The command above will look for the filename (`Test workflow.workflow.json`) in
the `/workflows` folder instead of the root folder.  All path values should be
relative to the root of the repository.

You can also specify the full path in the `filename` argument instead of using
the `--path` command line option.   For instance, the same command could also
be performed using the following syntax:

```
ipctl pull workflow "workflows/Test workflow.workflow.json" assets
```

Both forms of the command are acceptable.
