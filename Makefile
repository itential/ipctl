# !make

# Copyright 2024 Itential Inc. All Rights Reserved
# Unauthorized copying of this file, via any medium is strictly prohibited
# Proprietary and confidential

export GOOS        := $(shell uname | tr '[:upper:]' '[:lower:]')
export GOARCH      := amd64
export CGO_ENABLED := 0

.DEFAULT_GOAL := help

.PHONY: build \
	clean \
	config \
	coverage \
	licenses \
	snapshot \
	test 

# The help target displays a help message that includes the avialable targets
# in this `Makefile`.  It is the default target if `make` is run without any
# parameters.
help:
	@echo "Available targets:"
	@echo "  build      - Builds the iap application binary"
	@echo "  clean      - Cleans the development environment"
	@echo "  config     - Display the runtime config"
	@echo "  coverage   - Run test coverage report"
	@echo "  install    - Install application dependencies"
	@echo "  licenses   - Ensures that licenses exist on every code file and update license-attributions.md"
	@echo "  snapshot   - Create a development snapshot build"
	@echo "  test       - Run test suite"
	@echo ""

# The config target shows the current configured values that are enforce when
# any `make` target is run.
config: 
	@echo "GOOS=${GOOS}"
	@echo "GOARCH=${GOARCH}"
	@echo "CGO_ENABLED=${CGO_ENABLED}"
	@echo

# The clean target will remove all build and distribution directories as well
# as any coverage reports so the project directory is clean.
clean:
	@if [ -d bin ]; then rm -rf bin; fi
	@if [ -d dist ]; then rm -rf dist; fi
	@if [ -d cover ]; then rm -rf cover; fi

# The coverage target will run the unit test framework and generate a coverage
# report to see how much of the application source code has test coverage.
coverage:
	@scripts/test.sh coverage

# The licenses target is used to make sure that all source code files have the
# appropriate license header in them.   It will also generate a license file
# that includes all of the applicable licenses for any 3rd party libraries that
# have been included in this project.  To use this target, you must have the
# `go-licenses` tool installed and available in your path.  See
# `https://github.com/google/go-licenses` for more detaials
licenses:
	@go-licenses report . --template ./tools/license-attributions/template.tpl --ignore github.com/itential > license-attributions.md
	@go run ./tools/copyrighter/main.go

# The snapshot target will create a snapshot build of the application and place
# it into the dist/ folder.  The folder will be created if it doesn't already
# exist.  To use snapshot, you must have the `goreleaser` tool installed and
# available in your path.  See https://goreleaser.com for more details.
snapshot:
	BUILD=$$(git rev-parse --short HEAD) goreleaser release --snapshot --clean

# the install target will download the required go modules to the local cache
# and add and remove any missing or unused modules.
install:
	go mod download
	go mod tidy

# The build target will build the application binary and place it into the bin/
# folder.  If the folder does not exist, it will be created.
build: install
	go build \
		-v \
		-o bin/ipctl \
		-ldflags="-X 'github.com/itential/ipctl/internal/metadata.Build=$$(git rev-parse --short HEAD)' -X 'github.com/itential/ipctl/internal/metadata.Version=$$(git tag | sort -V | tail -1)'"

# The test target runs the unit tests for this project.   All unit tests are
# based on mock data for the various services it connects to.  The target
# should run without error before making any pull requests to this project.
test:
	@go run ./tools/copyrighter/main.go -check
	@scripts/test.sh unittest
