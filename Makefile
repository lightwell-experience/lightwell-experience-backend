GO_OUTPUT ?= release
BINARY := $(GO_OUTPUT)/lightwell-network
GOLANGCI_LINT := $(GO_OUTPUT)/golangci-lint
GOLANGCI_LINT_VERSION :=

GO_SOURCES := $(shell find cmd pkg -type f -name '*.go')

.PHONY: build
build: $(BINARY)

$(BINARY): $(GO_SOURCES) go.mod go.sum
	mkdir -p $(GO_OUTPUT)
	go build -o "$@" ./cmd/lightwell-network

.PHONY: run
run: $(BINARY)
	"$(BINARY)"

.PHONY: get-deps
get-deps:
	go get -d ./...

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	go test ./pkg/...

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT) ## Install golangci-lint into $(GO_OUTPUT)

$(GOLANGCI_LINT):
	mkdir -p $(GO_OUTPUT)
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(GO_OUTPUT) $(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Run Go linter
	$(GOLANGCI_LINT) run --fix --timeout=5m

.PHONY: clean
clean:
	rm -f "$(BINARY)"
