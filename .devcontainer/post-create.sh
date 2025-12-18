#!/bin/sh
set -e

git config --global --add safe.directory /workspaces/app

go mod tidy
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
