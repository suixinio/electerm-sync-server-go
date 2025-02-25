#!/bin/bash

# Make script exit on first error
set -e

echo "=== Running Electerm Sync Server Tests ==="
echo

# Create test data directory if it doesn't exist
mkdir -p test-data
cat > test.env << EOL
PORT=7837
HOST=127.0.0.1
JWT_SECRET=test-secret
JWT_USERS=testuser
FILE_STORE_PATH=./test-data
EOL

echo "Created temporary test environment"

# Install test dependencies if needed
go get github.com/stretchr/testify/assert

# Run the tests
go test -v ./src/...

# Clean up test data
rm -rf test-data
rm test.env

echo
echo "=== Tests Completed ==="