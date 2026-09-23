# JwtView: Development scripts

INSTALL_PATH := "$HOME/.local"

# Default command
_default:
    @just --list --unsorted

# Sync Go modules
tidy:
    go mod tidy
    @echo "All modules synced, Go workspace ready!"

# CLI local run wrapper
jwtview *args:
    @go run . {{ args }}

# Run all unit tests
test:
    #!/usr/bin/env bash
    echo "Running unit tests!"
    go clean -testcache
    go test -cover ./...

# Run all integration tests
integration-test:
    #!/usr/bin/env bash
    echo "Running integration tests!"
    go clean -testcache
    RUN_INTEGRATION=true go test -v -cover ./cmd/cli/...

# Build the binary
build:
    #!/usr/bin/env bash
    # Detect OS and architecture
    case "$(uname -s)" in
        Linux*) OS="linux" ;;
        Darwin*) OS="darwin" ;;
        *) echo "Error: Unsupported OS (${OS})"; exit 1 ;;
    esac
    case "$(uname -m)" in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        arm64) ARCH="arm64" ;;
        *) echo "Error: Unsupported architecture (${ENV_ARCH})"; exit 1 ;;
    esac

    echo "Building jwtview for ${OS}/${ARCH}..."
    go mod download all
    CGO_ENABLED=0 GOOS="${OS}" GOARCH="${ARCH}" go build -o ./jwtview .
    echo "Built binary for jwtview successfully!"

# Install the binary locally
install-local: build
    #!/usr/bin/env bash
    set -eux
    echo "Installing jwtview locally..."
    BIN_PATH="{{ INSTALL_PATH }}/bin/jwtview"
    cp ./jwtview "${BIN_PATH}"
    chmod +x "${BIN_PATH}"
    echo "Installed jwtview locally!"

# Remove the local binary
uninstall-local:
    #!/usr/bin/env bash
    set -eux
    echo "Uninstalling jwtview..."
    BIN_PATH="{{ INSTALL_PATH }}/bin/jwtview"
    rm "${BIN_PATH}"
    echo "Uninstalled jwtview!"

# Update the project dependencies
update-deps:
    @echo "Updating project dependencies..."
    go get -u ./...
    go mod tidy
