# Yor - Universal Infrastructure-as-Code Tagging
[![Maintained by Bridgecrew.io](https://img.shields.io/badge/maintained%20by-bridgecrew.io-blueviolet)](https://bridgecrew.io/?utm_source=github&utm_medium=organic_oss&utm_campaign=yor)
![golangci-lint](https://github.com/bridgecrewio/yor/workflows/tests/badge.svg)
[![security](https://github.com/bridgecrewio/yor/actions/workflows/security.yml/badge.svg)](https://github.com/bridgecrewio/yor/actions/workflows/security.yml)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-81%25-brightgreen.svg?longCache=true&style=flat)</a>
[![slack-community](https://slack.bridgecrew.io/badge.svg)](https://slack.bridgecrew.io/?utm_source=github&utm_medium=organic_oss&utm_campaign=yor)
[![Go Report Card](https://goreportcard.com/badge/github.com/bridgecrewio/yor)](https://goreportcard.com/report/github.com/bridgecrewio/yor)
[![Go Reference](https://pkg.go.dev/badge/github.com/bridgecrewio/yor.svg)](https://pkg.go.dev/github.com/bridgecrewio/yor)
 
Yor is an open-source tool that helps add informative and consistent tags across infrastructure-as-code frameworks such as Terraform, CloudFormation, and Serverless.

Yor is built to run as a [GitHub Action](https://github.com/bridgecrewio/yor-action) that hydrates IaC code with consistent tagging logics. It can also run as a pre-commit hook and a standalone CLI.

## Features
* Tracing: ```yor_trace``` tag enables simple attribution between an IaC resource block and a running cloud resource.
* Change management: git-based tags automatically add org, repo, commit and modifyer details on every resource block.  
* Custom taggers: user-defined tagging logics can be added to run using Yor.
* Skips: inline annotations enable developers to excluse paths that should not be tagged.

## Demo
### Attributing a resource to an owner
![](docs/yor_owner.gif)

### Change management tags
![](docs/yor_git_tags.gif)

### Trace IaC code to cloud resource
![](docs/yor_trace.gif)

### Trace cloud resource to IaC code
![](docs/yor_file.gif)

## **Table of contents**

- [Getting Started](#getting-started)
- [Disclaimer](#disclaimer)
- [Support](#support)
- [Customizing Yor](CUSTOMIZE.md)

## Getting Started

### Installation
GitHub Action
```yaml
- name: Checkout repo
  uses: actions/checkout@v2
  with:
    fetch-depth: 0
- name: Run yor action
  uses: bridgecrewio/yor-action@main
```

MacOS
```sh
brew tap bridgecrewio/tap
brew install bridgecrewio/tap/yor
```

Docker
```sh
docker pull bridgecrew/yor

docker run --tty --volume /local/path/to/tf:/tf bridgecrew/yor tag --directory /tf
```

Pre-commit
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

### Usage

`tag` : Apply tagging on a given directory.

```sh
 ./yor tag --directory terraform/
 # Apply all the tags in yor on the directory tree terraform.

 ./yor tag --directory terraform/ --skip-tags git_last_modified_by,yor_trace
 # Apply all the tags in yor except the tags git_last_modified_by and yor_trace.

 ./yor tag --tag-group git --directory terraform/
 # Apply only the tags under the git tag group.
```

`-o` : Modify output formats.

```sh
./yor tag -d . -o cli
# default cli output

./yor tag -d . -o json
# json output

./yor tag -d . --output cli --output-json-file result.json
# print cli output and additional output to file on json file -- enables prgormatic analysis alongside printing human readable result
```

`--skip-tags`:Specify only named tags (allow list) or run all tags except those listed (deny list).

```sh
./yor tag -d . --skip-tags yor_trace
## Run all but yor_trace

./yor tag -d . --skip-tags yor_trace,git_modifiers
## Run all but yor_trace and git_modifiers

./yor tag -d . --skip-tags git*
## Run all tags except tags with specified patterns
```

`skip-dirs` : Skip directoruy paths you can define paths that will not be tagged.

```sh
./yor tag -d path/to/files
## Run on the directory path/to/files

./yor tag -d path/to/files --skip-dirs path/to/files/skip,path/to/files/another/skip2
## Run yor on the directory path/to/files, skipping path/to/files/skip/ and path/to/files/another/skip2/
```

`list-tag`

```sh
./yor list-tag-groups
 # List tag classes that are built into yor.
 
 ./yor list-tags
 # List all the tags built into yor
./yor list-tags --tag-groups git
 
 # List all the tags built into yor under the tag group git
```

## Contributing

Contribution is welcomed! 

We are working on extending Yor and adding more parsers (to support additional IaC frameworks) and more taggers (to tag using other contextual data).

To maintain our conventions, please run lint on your branch before opening a PR. To run lint:
```sh
golangci-lint run --fix --skip-dirs tests/yor_plugins
```

## Disclaimer

`yor` does not save, publish or share with anyone any identifiable customer information.  
No identifiable customer information is used to query Bridgecrew's publicly accessible guides.

## Support

If you need direct support you can contact us at https://slack.bridgecrew.io/.
