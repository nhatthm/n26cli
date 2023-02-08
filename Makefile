CMD=n26

VENDOR_DIR = vendor
BUILD_DIR = build

GITHUB_OUTPUT ?= /dev/null

GOLANGCI_LINT_VERSION ?= v1.51.1

GO ?= go
GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint-$(GOLANGCI_LINT_VERSION)

$(info >> GOFLAGS: ${GOFLAGS})
ifneq "$(wildcard ./vendor )" ""
    $(info >> using vendor)
    modVendor =  -mod=vendor
    ifeq (,$(findstring -mod,$(GOFLAGS)))
        export GOFLAGS := ${GOFLAGS} ${modVendor}
    endif
endif

.PHONY: $(VENDOR_DIR) build build-linux lint test test-unit

$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)
	@$(GO) mod vendor
	@$(GO) mod tidy

lint: $(GOLANGCI_LINT) $(VENDOR_DIR)
	@$(GOLANGCI_LINT) run

test: test-unit test-integration

## Run unit tests
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

test-integration:
	@echo ">> integration test"
	@$(GO) test ./features/... -gcflags=-l -coverprofile=features.coverprofile -coverpkg ./... -godog -race -tags integration

build-linux:
	@echo ">> building binary, GOFLAGS: $(GOFLAGS)"
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(shell ./resources/scripts/build_args.sh)" -o $(BUILD_DIR)/$(CMD)-linux cmd/$(CMD)/*

## Build binary
build:
	@echo ">> building binary, GOFLAGS: $(GOFLAGS)"
	@rm -f $(BUILD_DIR)/*
	@$(GO) build -ldflags "$(shell ./resources/scripts/build_args.sh)" -o $(BUILD_DIR)/$(CMD) cmd/$(CMD)/*

.PHONY: $(GITHUB_OUTPUT)
$(GITHUB_OUTPUT):
	@echo "MODULE_NAME=$(MODULE_NAME)" >> "$@"
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)" >> "$@"

$(GOLANGCI_LINT):
	@echo "$(OK_COLOR)==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)$(NO_COLOR)"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin "$(GOLANGCI_LINT_VERSION)"
	@mv ./bin/golangci-lint $(GOLANGCI_LINT)
