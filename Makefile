# !make

# Copyright 2024 Itential Inc. All Rights Reserved
# GNU General Public License v3.0+ (see LICENSES/GPL-3.0-or-later.txt or https://www.gnu.org/licenses/gpl-3.0.txt)
# SPDX-License-Identifier: GPL-3.0-or-later

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

help:
	@echo "Available targets:"
	@echo "  build      - Builds the iap application binary"
	@echo "  clean      - Cleans the development environment"
	@echo "  config     - Display the runtime config"
	@echo "  coverage   - Run test coverage report"
	@echo "  licenses   - Ensures that licenses exist on every code file and update license-attributions.md"
	@echo "  snapshot   - Create a development snapshot build"
	@echo "  test       - Run test suite"
	@echo ""

config: 
	@echo "GOOS=${GOOS}"
	@echo "GOARCH=${GOARCH}"
	@echo "CGO_ENABLED=${CGO_ENABLED}"
	@echo

clean:
	@if [ -d bin ]; then rm -rf bin; fi
	@if [ -d dist ]; then rm -rf dist; fi
	@if [ -d cover ]; then rm -rf cover; fi

coverage:
	@scripts/test.sh coverage

licenses:
	@go-licenses report . --template ./tools/license-attributions/template.tpl --ignore github.com/itential > license-attributions.md
	@go run ./tools/copyrighter/main.go

snapshot:
	BUILD=$$(git rev-parse --short HEAD) goreleaser release --snapshot --clean

build:
	go build \
		-v \
		-o bin/ipctl \
		-ldflags="-X 'github.com/itential/ipctl/internal/metadata.Sha=$$(git rev-parse --short HEAD)' -X 'github.com/itential/ipctl/nternal/metadata.User=$$(id -u -n)' -X 'github.com/itential/ipctl/internal/metadata.Time=$$(date)' -X 'github.com/itential/ipctl/internal/metadata.Version=$$(git tag | sort -V | tail -1)'"

test:
	@scripts/test.sh unittest
