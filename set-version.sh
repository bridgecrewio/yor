new_tag=${{ needs.create-release.outputs.version }}
echo "new tag: $new_tag"
echo "package common" > src/common/version.go
echo "" >> src/common/version.go
echo "const Version = \"$new_tag\"" >> src/common/version.go