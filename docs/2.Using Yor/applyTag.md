---
layout: default
published: true
title: Applying Tags
nav_order: 3
---
# Applying Tags

The following commands are used to apply tags. GitHub Actions provides a simple, automatic way of applying tags to your IaC 
both during pull request review and as part of any build process. In order to integrate Yor into follow the installation 
[here](../2.Using Yor/installation.md#integrate-yor-with-github-actions)

## Apply Built-in Tags
To apply all configured tags run the following commands:

`./yor tag` - apply the built-in tags and any [custom](../3.Custom Taggers/customTagExamples.md) tags on a directory
   ```sh
    ./yor tag --directory terraform/
    # Apply all the tags in yor on the directory tree terraform/
   
    ./yor tag --directory terraform/ --skip-tags git_last_modified_by,yor_trace
    # Apply all the tags in yor except the tags git_last_modified_by and yor_trace
   
    ./yor tag --tag-group git --directory terraform/
    # Apply only the tags under the git tag group

   ```
## Tagging Docker Files

To run Yor as a Docker container, run the following commands after the file has been built.
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

## Use case: module tagging
Yor supports terraform [`module` blocks](https://www.terraform.io/docs/language/modules/sources.html) tagging using:
1. modules with a local path - will not be modified. The underlying resources will be tagged separately.
2. modules with a remote path - tags will be added according to the module block metadata.
   Yor does not download the remote module and modify it, but rather considers it as a black box.
   
Some examples:
```terraform
module "local_module" {
   # This is a local module. Yor will **not** modify this block. 
   # Instead, Yor will tag the actual resources located at the source dir that is specified in the module block
   source  = "../../tests/terraform"
   tags    = {
      env = var.env
   }
}

module "remote_module" {
   # This is a remote module (from the registry). 
   # Yor will add tags to the `tags` attribute of this module
   source = "terraform-aws-modules/vpc/aws"
   tags   = {
      env = var.env
   }
}

module "remote_module_2" {
   # This is a remote module (from github). 
   # Yor will add tags to the `tags` attribute of this module
   source = "git@github.com:terraform-aws-modules/terraform-aws-vpc.git"
   tags   = {
      env = var.env
   }
}
```

### Tagging examples:
#### Module with remote path
##### Before
```terraform
module "remote_module" {
   # This is a remote module (from the registry). 
   # Yor will add tags to the `tags` attribute of this module
   source = "terraform-aws-modules/vpc/aws"
   tags   = {
      env = var.env
   }
}
```

##### After
```terraform
module "remote_module" {
   # This is a remote module (from the registry). 
   # Yor will add tags to the `tags` attribute of this module
   source = "terraform-aws-modules/vpc/aws"
   tags   = {
      env                  = var.env
      yor_trace            = "912066a1-31a3-4a08-911b-0b06d9eac64e"
      git_repo             = "example"
      git_org              = "bridgecrewio"
      git_file             = "applyTag.md"
      git_commit           = "COMMITHASH"
      git_modifiers        = "bana/gandalf"
      git_last_modified_at = "2021-01-08 00:00:00"
      git_last_modified_by = "bana@bridgecrew.io"
   }
}
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
