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
```

## Integrate Yor with GitLab CI
Integrating Yor into a GitLab CI pipeline provides a simple, automatic way of applying tags to your IaC both
during pull request review and as part of any build process.
```yaml
.git-script: &git-script |
  cd $CI_PROJECT_DIR
  git status
  lines=$(git status -s | wc -l)
  if [ $lines -gt 0 ];then
    echo "committing"
    git config --global user.name "$AUTO_COMMITTER_NAME"
    git config --global user.email "$AUTO_COMMITTER_EMAIL"
    echo ".yor_plugins" >> .gitignore
    git add .
    git commit -m "YOR: Auto add/update yor.io tags."
    git push -o ci.skip "https://${GITLAB_USER_NAME}:${GIT_PUSH_TOKEN}@${CI_REPOSITORY_URL#*@}"
  else
    echo "no updated resources, nothing to commit."
  fi

stages:
  - yor

run-yor:    
  stage: yor
  script:
    - git checkout ${CI_COMMIT_REF_NAME}
    - export YOR_VERSION=0.1.62
    - wget -q -O - https://github.com/bridgecrewio/yor/releases/download/${YOR_VERSION}/yor-${YOR_VERSION}-linux-amd64.tar.gz | tar -xvz -C /tmp
    - /tmp/yor tag -d .
    - *git-script
```

You will need to set the following variables in `Settings > CI/CD > Variables` in order to use the pipeline.
`AUTO_COMMITTER_EMAIL`: Configure the author's e-mail for the git commit of new tags.
`AUTO_COMMITTER_NAME`: Configure the authors name for the git commit of new tags.
`GITLAB_USER_NAME`: The GitLab username for authenticating the git push of new tags, must match a valid users `GIT_PUSH_TOKEN`.
`GIT_PUSH_TOKEN`: A GitLab personal access token with permissions to commit to the repository.


## MacOS
Run the following commands to install Yor on MacOS with [Homebrew](https://brew.sh/):
```sh
brew tap bridgecrewio/tap
brew install bridgecrewio/tap/yor
```

## Windows
Run the following command to install Yor on Windows with [Chocolatey](https://chocolatey.org/install):
```sh
choco install yor
```

## Docker
To install and run Yor using a Docker image, run the following commands after the file has been built.
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
