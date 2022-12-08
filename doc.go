//go:build tools
// +build tools

package main

import (
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "mvdan.cc/gofumpt"
)
