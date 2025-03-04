//go:build tools
// +build tools

package tools

import (
	_ "github.com/derision-test/go-mockgen/cmd/go-mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/vektra/mockery/v2"
	_ "gotest.tools/gotestsum"
	_ "mvdan.cc/gofumpt"
)
