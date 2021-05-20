---
layout: default
published: true
title: Installing Yor
nav_order: 1
---

# Installing Yor
You can install Yor in a number of different ways:

## Integrate Yor with Github Actions
Integrating Yor into GitHub Actions provides a simple, automatic way of applying tags to your IaC both
during pull request review and as part of any build process.
```yaml
- name: Checkout repo
  uses: actions/checkout@v2
  with:
    fetch-depth: 0
- name: Run yor action
  uses: bridgecrewio/yor-action@main
- name: Commit tag changes
  uses: stefanzweifel/git-auto-commit-action@v4
```

## MacOS
Run the following commands to install Yor on MacOS:
```sh
brew tap bridgecrewio/tap
brew install bridgecrewio/tap/yor
```
```sh
brew install yor
```

## Docker
To install Yor and add tags to your Dockerfile, run the following commands after the file has been built.
```sh
docker pull bridgecrew/yor

docker run --tty --volume /local/path/to/tf:/tf bridgecrew/yor tag --directory /tf
```

## Pre-commit
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