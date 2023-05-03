//go:build tools
// +build tools

package tools

import (
	_ "github.com/derision-test/go-mockgen/cmd/go-mockgen"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "mvdan.cc/gofumpt"
)
