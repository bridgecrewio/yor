# Yor
![golangci-lint](https://github.com/bridgecrewio/yor/workflows/tests/badge.svg)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-76%25-brightgreen.svg?longCache=true&style=flat)</a>

Automated IaC tagging using external sources of data.

Yor applies IaC tags based on git data so that you can consistently attribute changes to cloud resources managed using infrastructure-as-code. Yor is currently built to tag based on local or remote Git data. It is built 



## Features

* Yor collects data from [git-blame](https://git-scm.com/docs/git-blame) and enables mapping individual resources to specific commits. It can be run automatically  
* ```yor_trace``` tag enables full attribution between build time and run time resources. 
* ```git_*``` tags  connect cloud resources to individual git commits and enable assigning clear ownership between developers and the resources they routinely change.

 

### Supported tags

```
yor_trace = "842fa55e-c0f8-4f79-a56d-a046a24d8e08"
git_org = "bridgecrewio"
git_repo = "terragoat"
git_file = "README.md" # this is the path from the repo root dir...
git_commit = "47accf06f13b503f3bab06fed7860e72f7523cac" # This is the latest commit for this resource
git_last_modified_At = "2020-03-28 21:42:46 +0000 UTC"
git_last_modified_by = "schosterbarak@gmail.com"
git_modifiers = "schosterbarak/baraks" # These are extracted from the emails, everything before the @ sign. Can be done for modified_by tag as well
```



## **Table of contents**

- [Getting Started](#getting-started)
- [Disclaimer](#disclaimer)
- [Support](#support)

## Getting Started

### Installation

On MacOS

```sh
brew install yor
```



### Usage

* Tag an entire directory

```sh
yor --directory terraform/
```

### Skipping tags 

Using command line flags you can specify to run only named tags (allow list) or run all tags except 
those listed (deny list).

```sh
yor -d . --tag yor_trace
## Run only yor_trace

yor -d . --skip-tag yor_trace
## Run all but yor_trace

yor -d . --skip-tag git*
## Run all tags except tags with specified patterns

yor -d . --skip-tag
```

### Output formats

```sh
yor -d . -o cli
# default cli output

yor -d . -o json
# json output

yor -d . --output cli --output-json-file result.json
# will print cli output and additional output to file on json file -- enables prgormatic analysis alongside printing human readable result
```

### Loading External Tags Using Plugins

#### Prerequisites

An example can be found in `tests/yor_plugins/example`

1. Create tags implementing the `ITag` interface.
2. If you wish to override an existing tag, make the tag's method `GetPriority()` return a positive number. Otherwise, return `0` or a negative number.
3. Create a file located in package `main` that exposes a variable `ExtraTags` - array containing pointers to all tags implemented.
4. Run command `go build -gcflags="all=-N -l" -buildmode=plugin -o <plugin-dir>/extra_tags.so <plugin-dir>/*.go`

```go
package main

var ExtraTags = []interface{}{&TerragoatTag{}, &CheckovTag{}}
```

#### Running yor

```sh
yor --custom-tagger tests/yor_plugins/example
# run yor with custom tags located in tests/yor_plugins/example
```


#### Troubleshooting
If you encounter the following error: 
`plugin was built with a different version of package ...`

Build yor with `go build -gcflags="all=-N -l"`


## Contributing

Contribution is welcomed! 

We are working on extending Yor and adding more parsers (to support additional IaC frameworks) and more taggers (to tag using other contextual data).

## Disclaimer

`yor` does not save, publish or share with anyone any identifiable customer information.  
No identifiable customer information is used to query Bridgecrew's publicly accessible guides.

## Support

If you need direct support you can contact us at info@bridgecrew.io.
