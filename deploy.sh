echo "Fetching the latest code from repo" 
git pull

echo "Setting go version in go.mod to 1.23"

#!/bin/bash

# Check if go.mod file exists
if [[ ! -f go.mod ]]; then
  echo "go.mod file not found."
  exit 1
fi

# Read the current go.mod content
current_content=$(cat go.mod)

# Replace the version string
new_content=$(echo "$current_content" | sed 's/go 1.\([0-9]\)/go 1.23/g')

# Check if the content has changed
if [[ "$new_content" != "$current_content" ]]; then
  # Write the new content to go.mod
  echo "$new_content" > go.mod
  echo "Updated go.mod to version 1.23"
else
  echo "go.mod version is already 1.23"
fi


echo "Compiling..."
go build -o sg-biz-service

echo "Removed older binary"
rm ../sg-biz-service

echo "New binary moved to parent"
cp sg-biz-service ../

echo "Restarting the service"
systemctl restart sg-biz-service
