# Yor - Infrastructure as code tagging engine
[![Maintained by Bridgecrew.io](https://img.shields.io/badge/maintained%20by-bridgecrew.io-blueviolet)](https://bridgecrew.io/?utm_source=github&utm_medium=organic_oss&utm_campaign=yor)
![golangci-lint](https://github.com/bridgecrewio/yor/workflows/tests/badge.svg)
[![security](https://github.com/bridgecrewio/yor/actions/workflows/security.yml/badge.svg)](https://github.com/bridgecrewio/yor/actions/workflows/security.yml)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-82%25-brightgreen.svg?longCache=true&style=flat)</a>
[![slack-community](https://slack.bridgecrew.io/badge.svg)](https://slack.bridgecrew.io/?utm_source=github&utm_medium=organic_oss&utm_campaign=yor)
[![Go Report Card](https://goreportcard.com/badge/github.com/bridgecrewio/yor)](https://goreportcard.com/report/github.com/bridgecrewio/yor)
[![Go Reference](https://pkg.go.dev/badge/github.com/bridgecrewio/yor.svg)](https://pkg.go.dev/github.com/bridgecrewio/yor)
 
Cloud service providers allow users to assign metadata to their cloud resources in the form
of tags. Each tag is a simple label consisting of a customer-defined key and a value
that can make it easier to manage, search for, and filter resources. Although there are no
inherent types of tags, they enable customers to categorize resources by purpose, owner,
environment, or other criteria.

Tags can be used for security, cost allocation, automation, console organization, access control, and operations. 

Yor is an open-source tool that helps to manage tags in a consistent manner across infrastructure as code frameworks (Terraform, Cloudformation, and Serverless Framework) .
By auto-tagging in IaC you will be able to trace any cloud resource from code to cloud. 
Yor enables version-controlled owner assignment and resource tracing based git history. It also can extend tag enforcement logic by loading external tagging logic into the CI/CD pipeline. 

## Features
* Simple tagging - tag a directory IaC resource by user input
* Git tagging - Yor collects data from [git-blame](https://git-scm.com/docs/git-blame) and enables mapping individual resources to specific commits and users.
* ```yor_trace``` tag enables full attribution between build time and run time resources.

## Demo
> Generate an git tags on infrastructure as code
![](docs/yor_git_tags.gif)

> Track resource last modifying user from code to cloud
![](docs/yor_owner.gif)

> Track resource identifer from code to cloud
![](docs/yor_trace.gif)

> Track resource code block file cloud to code
![](docs/yor_file.gif)

## **Table of contents**

- [Getting Started](#getting-started)
- [Disclaimer](#disclaimer)
- [Support](#support)

More Docs:
- [Customizing Yor](CUSTOMIZE.md)

## Getting Started

### Installation

On MacOS

```sh
brew tap bridgecrewio/tap
brew install bridgecrewio/tap/yor
```



### Usage

Yor supports the following commands:
1. `list-tag-groups` - list the groups of tags that are built into yor
   ```sh
   ./yor list-tag-groups
   ```
2. `list-tags` - list all the tags yor has built in. This will print each tag key
and the relevant group the tag belongs to.
   ```sh
    ./yor list-tags
    # List all the tags built into yor
   ./yor list-tags --tag-groups git
    # List all the tags built into yor under the tag group git
    ```
3. `tag` - apply the built in tags and any [custom](CUSTOMIZE.md) tags on a directory
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

### Output formats

```sh
./yor tag -d . -o cli
# default cli output

./yor tag -d . -o json
# json output

./yor tag -d . --output cli --output-json-file result.json
# will print cli output and additional output to file on json file -- enables prgormatic analysis alongside printing human readable result
```


#### Troubleshooting
If you encounter the following error: 
`plugin was built with a different version of package ...`

Build yor with `go build -gcflags="all=-N -l"`


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
