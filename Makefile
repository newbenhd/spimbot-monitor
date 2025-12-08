BINARY_NAME=spimbot-monitor

# Build directories
BUILD_DIR=build

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Default target builds both platforms
all: build

# Build for both darwin/arm64 and linux/amd64
build: build-darwin build-linux
	@echo "Build complete for both architectures"

# Build for darwin/arm64 (macOS Apple Silicon)
build-darwin:
	@echo "Building for darwin/arm64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

# Build for linux/amd64
build-linux:
	@echo "Building for linux/amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
