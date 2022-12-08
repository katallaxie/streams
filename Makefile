.DEFAULT_GOAL := build

GO_TEST = go run gotest.tools/gotestsum --format pkgname

.PHONY: generate
generate:
	go generate ./...

.PHONY: build ## Build the binary file.
build: generate fmt vet lint test
	GOPRIVATE=github.com/ionos-cloud go mod tidy
	goreleaser build --snapshot --rm-dist

.PHONY: fmt
fmt: ## Run go fmt against code.
	go run mvdan.cc/gofumpt -w .

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: generate fmt ## Run tests.
	mkdir -p .test/reports
	$(GO_TEST) --junitfile .test/reports/unit-test.xml -- -race ./... -count=1 -short -cover -coverprofile .test/reports/unit-test-coverage.out

.PHONY: lint
lint: ## Run golangci-lint against code.
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 2m

.PHONY: clean
clean: ## Remove previous build.
	find . -type f -name '*.gen.go' -exec rm {} +
	git checkout go.mod
