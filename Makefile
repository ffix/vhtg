APP_NAME := vhtg
GO_CMD := go
GO_BUILD := $(GO_CMD) build
GO_PARAMS := ./cmd


$(APP_NAME): cmd/main.go
	$(GO_BUILD) -o $(APP_NAME) $(GO_PARAMS)

# macOS (AMD64)
$(APP_NAME)-darwin-amd64: cmd/main.go
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o $(APP_NAME)-darwin-amd64 $(GO_PARAMS)

# macOS (M1 ARM64)
$(APP_NAME)-darwin-arm64: cmd/main.go
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) -o $(APP_NAME)-darwin-arm64 $(GO_PARAMS)

# Linux (AMD64)
$(APP_NAME)-linux-amd64: cmd/main.go
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(APP_NAME)-linux-amd64 $(GO_PARAMS)

# Linux (AMD64, static)
$(APP_NAME)-linux-amd64-static: cmd/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -a -ldflags '-extldflags "-static"' -o $(APP_NAME)-linux-amd64-static $(GO_PARAMS)

# Linux (ARM64)
$(APP_NAME)-linux-arm64: cmd/main.go
	GOOS=linux GOARCH=arm64 $(GO_BUILD) -o $(APP_NAME)-linux-arm64 $(GO_PARAMS)

# Build all targets
all: $(APP_NAME)-darwin-amd64 $(APP_NAME)-darwin-arm64 $(APP_NAME)-linux-amd64 $(APP_NAME)-linux-arm64

# Clean up the built binaries
clean:
	rm -f $(APP_NAME)-darwin-amd64 $(APP_NAME)-darwin-arm64 $(APP_NAME)-linux-amd64 $(APP_NAME)-linux-arm64

.PHONY: all clean
