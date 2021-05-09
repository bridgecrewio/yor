# Applying Tags

The following commands are used to apply tags.
`./yor` 
`tag` - apply the built in tags and any [custom](CUSTOMIZE.md) tags on a directory
   ```sh
    ./yor tag --directory terraform/
    # Apply all the tags in yor on the directory tree terraform/
   
    ./yor tag --directory terraform/ --skip-tags git_last_modified_by,yor_trace
    # Apply all the tags in yor except the tags git_last_modified_by and yor_trace
   
    ./yor tag --tag-group git --directory terraform/
    # Apply only the tags under the git tag group
    ```

Using docker:
```sh
docker pull bridgecrew/yor

docker run --tty --volume /local/path/to/tf:/tf bridgecrew/yor tag --directory /tf
```

Using pre-commit:

Add a hook:

```yaml
  - repo: git://github.com/bridgecrewio/yor
    rev: 0.0.44
    hooks:
      - id: yor
        name: yor
        entry: yor tag -d
        args: ["example/examplea"]
        language: golang
        types: [terraform]
        pass_filenames: false
```

To your **.pre-commit-config.yaml** and change the args and version number.

### Skipping tags

Using command line flags you can specify to run only named tags (allow list) or run all tags except
those listed (deny list).

```sh
./yor tag -d . --skip-tags yor_trace
## Run all but yor_trace

./yor tag -d . --skip-tags yor_trace,git_modifiers
## Run all but yor_trace and git_modifiers

./yor tag -d . --skip-tags git*
## Run all tags except tags with specified patterns
```

### Skipping directories

Using the command line flag skip-paths you can define paths which won't be tagged.
Be mindful that the skipped path should include the root dir path. See example below:

```sh
./yor tag -d path/to/files
## Run on the directory path/to/files

./yor tag -d path/to/files --skip-dirs path/to/files/skip,path/to/files/another/skip2
## Run yor on the directory path/to/files, skipping path/to/files/skip/ and path/to/files/another/skip2/
```