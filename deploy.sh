#!/bin/bash

echo "Stashing the local changes"
git stash

echo "Fetching the latest code from repo" 
git pull

echo "Setting go version in go.mod to 1.23"

# Check if go.mod exists
if [ ! -f go.mod ]; then
    echo "Error: go.mod file not found in the current directory."
    exit 1
fi

# Use sed to replace the Go version
sed -i 's/^go \(1\.[0-9]\+\(\.[0-9]\+\)\?\)/go 1.23/' go.mod

# Check if the replacement was successful
if grep -q "^go 1.23$" go.mod; then
    echo "Go version successfully updated to 1.23 in go.mod"
else
    echo "Failed to update Go version. Please check the go.mod file manually."
fi

echo "Compiling..."
go build -o sg-biz-service

echo "Removed older binary"
rm ../sg-biz-service

echo "New binary moved to parent"
cp sg-biz-service ../

echo "Restarting the service"
systemctl restart sg-biz.service
