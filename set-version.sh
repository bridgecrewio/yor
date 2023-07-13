#!/bin/sh
set -e
latestTag=$(git describe --tags)
echo "Updating version file with new tag: $latestTag"
echo "package common" > src/common/version.go
echo "" >> src/common/version.go
echo "const Version = \"$latestTag\"" >> src/common/version.go
