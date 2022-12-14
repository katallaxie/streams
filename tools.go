//go:build tools
// +build tools

package streams

import (
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "gotest.tools/gotestsum"
	_ "mvdan.cc/gofumpt"
)
