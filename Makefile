.PHONY: build install clean test build-all

VERSION := 1.0.0

# Default: build for current platform
build:
	@echo "🔨 Building prevent-llm-delete for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@go build -ldflags="-s -w" -o prevent-llm-delete main.go
	@echo "✅ Built: prevent-llm-delete"

# Build for all platforms
build-all:
	@echo "🔨 Building for all platforms..."
	@mkdir -p dist

	@echo "Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/prevent-llm-delete-darwin-amd64 main.go

	@echo "Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/prevent-llm-delete-darwin-arm64 main.go

	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/prevent-llm-delete-linux-amd64 main.go

	@echo "Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/prevent-llm-delete-linux-arm64 main.go

	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/prevent-llm-delete-windows-amd64.exe main.go

	@echo "✅ All binaries built in dist/"

# Install locally
install: build
	@echo "📦 Installing prevent-llm-delete to /usr/local/bin..."
	@sudo mv prevent-llm-delete /usr/local/bin/
	@echo "✅ Installed! Run 'prevent-llm-delete --help' to get started"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -f prevent-llm-delete
	@rm -rf dist/
	@echo "✅ Clean complete"

# Test the binary
test: build
	@echo "🧪 Testing prevent-llm-delete..."
	@./prevent-llm-delete status
	@echo "✅ Tests passed"

# Create release archives
release: build-all
	@echo "📦 Creating release archives..."
	@cd dist && tar -czf prevent-llm-delete-darwin-amd64.tar.gz prevent-llm-delete-darwin-amd64
	@cd dist && tar -czf prevent-llm-delete-darwin-arm64.tar.gz prevent-llm-delete-darwin-arm64
	@cd dist && tar -czf prevent-llm-delete-linux-amd64.tar.gz prevent-llm-delete-linux-amd64
	@cd dist && tar -czf prevent-llm-delete-linux-arm64.tar.gz prevent-llm-delete-linux-arm64
	@cd dist && zip prevent-llm-delete-windows-amd64.zip prevent-llm-delete-windows-amd64.exe
	@echo "✅ Release archives created in dist/"
