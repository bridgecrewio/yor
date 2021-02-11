# Yor
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-77%25-brightgreen.svg?longCache=true&style=flat)</a>

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

## Contributing

Contribution is welcomed! 

We are working on extending Yor and adding more parsers (to support additional IaC frameworks) and more taggers (to tag using other contextual data).

## Disclaimer

`yor` does not save, publish or share with anyone any identifiable customer information.  
No identifiable customer information is used to query Bridgecrew's publicly accessible guides.

## Support

If you need direct support you can contact us at info@bridgecrew.io.
