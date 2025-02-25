#!/bin/bash

echo "=== Formatting Go Code ==="

# Format all Go files
go fmt ./...

# For more aggressive formatting (handles imports etc.)
echo "Running goimports..."
if ! command -v goimports &> /dev/null; then
    echo "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
    # Use the binary directly from GOPATH/bin
    GOPATH=$(go env GOPATH)
    GOIMPORTS="$GOPATH/bin/goimports"
else
    GOIMPORTS="goimports"
fi

$GOIMPORTS -w ./src

echo "=== Formatting Complete ==="