#!/bin/bash
set -e

echo "Running unit tests..."
go test ./... -v -race -cover

echo "All tests passed ✅"