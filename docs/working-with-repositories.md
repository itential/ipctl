# Working with repositories

One of the key features of `ipctl` is the opportunity to work natively with
`git` repositories.  This integration allows you to directly push and pull
Itential Platform assets to and from git repositories using the `import` and
`export` commands with command line options.

## Git configuration

The `ipctl` application provides the following configuration options when
working with Git repositories.  All Git repository configuration options are
configured in the `[git]` configuration section.

- `name`
- `email`

The `name` configuration setting allows you to set the name of the user to use
when making commits using the `export` function.  The default value is the
current logged username.

The `email` configuration setting allows you to configure the email address of
the user to use when making commits using the `export` function.

For example, the following example demonstrates how to configure Git support
for a user called `Itential Pronghorn`.

```ini
[git]
name = Itential Pronghorn
email = pronghorn@itential.com
```

Once these values are set, they are used when pushing commits to Git
repositories using the `export` command.

## Operation

The `import` and `export` commands have been updated to work directly with Git
repositories instead of only the local file system.  In order to import and
export assets directly to and from Git repositories, the following command line
options have been added to `ipctl`.

- `--repository`
- `--reference`
- `--private-key-file`
- `--message`
- `--path`

The `--repository` command line option is required when pushing to a Git
repository.  This command line option accepts any valid Git URL which is used
to push the exported document to.

The `--reference` command line option is optional and is used to specify the Git
reference to push the asset to.  The reference can be a Git branch name, a Git
branch tag or a specific Git SHA.

The `--private-key-file` command line option is optional and defines the path
to the private key file to use when connecting to the Git repository.  This
option is used in conjunction with Git repositories that are connected to over
SSH.  If this option is not specified, the `ipctl` application will default to
using the local private key file.

The `--message` command line option is optional and sets the commit message to
include when committing the asset to the Git repository.  If this command line
option is not specified, a default commit message will be used.

The `--path` command line option is optionally used to define the path relative
to the root of the Git repository to store the asset.  When the `--repository`
command line option is not specified, this option will perform the same
function relative to the local current working directory.

## Pushing assets to repositories using `export`

The `export` command can be used to export an asset from an Itential Platform
server and push it directly into a Git repository.  The repository must already
exist and be initialized for this feature to work.

The `export` command uses the following command line options for working
with Git repositories:

- `--repository`
- `--reference`
- `--private-key-file`
- `--message`
- `--path`

Below is a very basic example of exporting an asset directly to a Git
repository.

```bash
ipctl export project "My Test Project" \
    --repository git@github.com:itential/assets.git \
    --message "committing new project to repository"
```

The above example command will export the project named `My Test Project` and
push it into a Git repository.  The file will be stored in the root of the
repository.

In some cases, it may be desirable to store the file in a specific folder
within the repository.  In order to do that, use the `--path` command line
option.

For example, assume the same project as before needs to be stored in folder
called `test/projects`.  The command would be

```bash
ipctl export project "My Test Project" \
    --repository git@github.com:itential/assets.git \
    --message "committing new project to repository" \
    --path test/projects
```

The added path is always relative to the root of the Git repository.  If the
path does not exist in the repository, it will be created.

## Pulling assets from repositories using `import`

The `import` command can be used to pull an asset directly from a Git
repository and import it into an Itential Platform server.  The `import`
command accepts the following command line options for performing this task:

- `--repository`
- `--reference`

Using these command line options will instruct `ipctl` to pull the asset from
the specified Git repository and import it directly to an Itential Platform
server.

For example, below is a very basic example that will import a project file
directly from a Git repository.

```bash
ipctl import project "My Test Project.project.json" \
    --repository git@github.com:itential/assets.git
```

When run, the command will connect to the Git repository defined by
`--repository` and extract the project file `My\ Test\ Project.project.json`
which will then be imported to the Itential Platform server.

The `import` command accepts a relative path as an argument to specify a
project file that is nested in a folder.  For example, the following command
will import the project file from the `test/projects` folder of the Git
repository.

```bash
ipctl import project "test/projects/My Test Project.project.json" \
    --repository git@github.com:itential/assets.git
```
