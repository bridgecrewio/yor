---
layout: default
published: true
title: Applying Tags
nav_order: 3
---
# Applying Tags

The following commands are used to apply tags.

## Apply Built-in Tags
To apply all configured tags run the following commands:

`./yor tag` - apply the built-in tags and any [custom](/docs/3.Custom Taggers/customTagExamples.md) tags on a directory
   ```sh
    ./yor tag --directory terraform/
    # Apply all the tags in yor on the directory tree terraform/
   
    ./yor tag --directory terraform/ --skip-tags git_last_modified_by,yor_trace
    # Apply all the tags in yor except the tags git_last_modified_by and yor_trace
   
    ./yor tag --tag-group git --directory terraform/
    # Apply only the tags under the git tag group

   ```
## Tagging Docker Files

To add tags to your Dockerfile, run the following commands after the file has been built.
```sh
docker pull bridgecrew/yor

docker run --tty --volume /local/path/to/tf:/tf bridgecrew/yor tag --directory /tf
```

## Tagging Using Pre-commit:
Using Pre-commit with Yor provides a simple, automatic way of applying tags to your IaC identifying potential issues before submission to code review.

You need to have the pre-commit package manager installed before you can run Pre-commit hooks.

Add a hook to your **.pre-commit-config.yaml** and change the args and version number.

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

## Skipping Tags

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

## Skipping Directories

Using the command line flag skip-paths you can define paths which won't be tagged.
Be mindful that the skipped path should include the root dir path. See example below:

```sh
./yor tag -d path/to/files
## Run on the directory path/to/files

./yor tag -d path/to/files --skip-dirs path/to/files/skip,path/to/files/another/skip2
## Run yor on the directory path/to/files, skipping path/to/files/skip/ and path/to/files/another/skip2/

```
