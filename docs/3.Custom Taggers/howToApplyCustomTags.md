# Applying Custom Tags

## Running yor with the external tags / taggers

```sh
./yor tag --custom-tagging tests/yor_plugins/example
# run yor with custom tags located in tests/yor_plugins/example

./yor tag --custom-tagging tests/yor_plugins/example,tests/yor_plugins/tag_group_example
# run yor with custom tags located in tests/yor_plugins/example and custom taggers located in tests/yor_plugins/tag_group_example
```