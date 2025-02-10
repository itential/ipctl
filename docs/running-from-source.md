# Running from source

The `ipctl` application can be direclty run from a git checkout.  In order to
run the application directly from the source tree, you must have `golang`
installed  and available on in your local environment.

To get started, first install `golang` 1.22 on your local system.  Installation
instruction can be found [here](https://go.dev/doc/install).

# Running from source

The repository includes a `Makefile` with targets to allow installing and
building the application.  Be sure to install `make` into your local
environment as well.

Once installed and from the root directory where you have cloned the
repository, and run `make install`.   This will install of the `golang`
dependencies and prepare your system to run the code from source.

Finally to execute the code from source run `go run main.go` which will run
the application directly from source code.

## Building the executable

You can also build the application executable using the `build` make target. To
build the application simply run `make build` from the root folder of the
repository.   This will build the executable and place it into the `bin/`
folder.

## Building snapshots

Finally, the `Makefile` also supports building snapshots.  Snapshots will build
the binary for all supported platforms.   To create a snapshot build, you will
need to have `goreleaser` installed and available in your local environment.
See [here](https://goreleaser.com/) for details on how to install `goreleaser`.

Once `goreleaser` is installed, you can build your own snapshot using `make
snapshot`  The comand will build the application for all supported platforms
and place them into the `dist/` folder.
