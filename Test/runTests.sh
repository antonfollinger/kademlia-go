#!/bin/bash

COVERAGE_FILE="coverage.out"
TARGET_DIR="./internal/kademlia"

cd ..

go test -coverprofile="$COVERAGE_FILE" "$TARGET_DIR"

# Display coverage summary
go tool cover -func="$COVERAGE_FILE"

# Generate an HTML coverage report
go tool cover -html="$COVERAGE_FILE" -o coverage.html

# Move coverage files to the Test directory
mv "$COVERAGE_FILE" ./Test/
mv coverage.html ./Test/