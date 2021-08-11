#!/bin/sh
echo "Updating version file with new tag: $GORELEASER_CURRENT_TAG"
echo "package common" > src/common/version.go
echo "" >> src/common/version.go
echo "const Version = \"$GORELEASER_CURRENT_TAG\"" >> src/common/version.go