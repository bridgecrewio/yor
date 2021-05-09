# Installing Yor

GitHub Action
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

MacOS
```sh
brew tap bridgecrewio/tap
brew install bridgecrewio/tap/yor
```
```sh
brew install yor
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