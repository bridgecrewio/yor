---
layout: default
published: true
title: Custom Tagger Using YAML  Configuration Files
nav_order: 2
---
# Custom Tagger Using YAML Configuration Files

The Yor framework uses YAML configuration files to support advanced rules when applying custom tags.
Users can define tagging enforcement rules that are specific to their organization’s needs.
YAML based custom tagging enables you to have different tags for different existing resource tags.

## Running YAML based custom tagger
In the CLI, define the path of the YAML configuration file that you want to apply. For example:

`yor tag -d . --config-file </path/to/file>`

You can also decide to use code owners file in case of git modifiers has a conflict (default: false)

`yor tag -d . --config-file </path/to/file> --use-code-owners`

The YAML based custom tagging configuration file includes the following options:
1. Name key definition (required)
2. Default key value and default keys (required)

***Example 1:** basic key-value tagging in all IaC resources*

```
name: env
value:
    default: prod
```

3. Filter definition (optional) - use cases where tagging will be applied:
    1. *Tags sequence*: tag resources that have all the tags within the map
4. Directory definition (optional): path to defined taggable resources

***Example 2:** Tagging specific resource types with specific `key:value` tags in a defined directory.*
```
name: env
value:
    default: prod
filters:
    tags:
        git_modifiers: donnaj
        git_repo: checkov
    directory: /path/to/some/dir
```

5. Use case dynamic value definition using *value* mapping (optional): Tags are defined based on matching
   keys that contain a sequence of values. Under each value the user can define which existing tags a resource will be
   tagged with. If none of the conditions are matched, a default value will be applied. In the example below
   resources in the directory `/path/to/some/dir` and existing tag `yor_trace: 123` will
   be tagged with one of the following:
    1. *team: devops*: resources have the tags `git_repo: yor`, `git_commit: asd12f`, and `git_modifiers:`
       will be tagged with one of the following values - `johnb / amyh / rond`
    2. *team: dev1*: for any other resource that does not comply with option 1.

```
name: team
value:
   default: dev1
   matches:
        - devops:
              tags:
                   git_modifiers:
                       - johnb
                       - amyh
                       - rond
                   git_commit: asd12f
                   git_repo: yor
   tags:
        yor_trace: 123
   directory: /path/to/some/dir
```
6. You can create custom tag groups using the YAML-based configuration. Add the `tag_groups` field, add a
   `name`, and a `tags` sequence which includes the tag items mentioned above. In the example below, shows a tag group
   named `ownership` which includes the two custom tags: `env` and `team`.

```
tag_groups:
  - name: ownership
    tags:
      - name: env
        value:
          default: dev
        filters:
          tags:
            git_repo: yor
            git_modifiers: tronxd
      - name: team
        value:
          default: interfaces
          matches:
            - seceng:
                tags:
                  git_modifiers:
                    - rotemavni
                    - tronxd
                    - nimrodkor
            - platform:
                tags:
                  git_modifiers:
                    - milkana
                    - nofar
        filters:
          tags:
            git_commit: 00193660c248483862c06e2ae96111adfcb683af
```

## Custom tagging using CLI

Some YAML configuration capabilities are available in the CLI. Some commands available are:
1. `--tag-name`: define tag name
2. `--tag-value`: define tag value
3. `-filter-tags`: tag resources that have tags as defined. Use an array [] to support multiple values and to support `AND` logic between tags

In the example below, EC2 instances and Security Groups will be tagged with the `env:prod` tag. Use this in cases where a resource that has `tronxd`
or `amy` are one of the `git_modifiers` and it is located in `checkov` or `terragoat git_repo`.

**Example 3:** CLI custom tagging

```sh
yor tag --tag-name env –tag-value prod –filter-tags git_modifiers=[tronxd,amy];git_repo=[checkov,terragoat]
```

## Running Yor with Custom Taggers
Use the following example to run Yor with custom tags:
```sh
./yor tag --custom-tagging tests/yor_plugins/example
# run yor with custom tags located in tests/yor_plugins/example

./yor tag --custom-tagging tests/yor_plugins/example,tests/yor_plugins/tag_group_example
# run yor with custom tags located in tests/yor_plugins/example and custom taggers located in tests/yor_plugins/tag_group_example
```

## Security considerations for `--config-file`

The file you pass to `--config-file` is **trusted input**: Yor reads tag
names, default values, and match values from it and writes them directly into
your IaC files. Treat the config file with the same care you treat your
CI/CD pipeline definitions.

### Do not point `--config-file` at an untrusted or repo-resident path

Storing the config file inside the same directory you scan with `-d` (a
common pattern — the config lives next to the IaC it tags) is convenient,
but it means **anyone who can open a pull request can edit the tagging
rules**. In particular, the config supports a `${env:NAME}` template syntax
that Yor expands by reading the named environment variable from the runner
and writing the expanded value as a tag. If the runner has CI secrets in
scope (the standard GitHub Action setup, including the
`pull_request_target` trigger), an attacker can reference a secret variable
from the config and have Yor copy that secret value into your IaC files —
which the [official `bridgecrewio/yor-action`](https://github.com/bridgecrewio/yor-action)
then commits and pushes back to the repository, persisting the leak to git
history.

**Recommended layouts:**

| Layout | Safe? |
| --- | --- |
| `--config-file` lives **outside** the directory passed to `-d` (e.g. checked into a separate, protected repo or assembled from CI variables at runtime) | ✅ Yes |
| `--config-file` lives **inside** `-d`, but the action runs with `commit_changes: 'NO'` (Yor only reports, never pushes) | ✅ Yes |
| `--config-file` lives **inside** `-d` and the action auto-commits | ❌ **Avoid** — a PR author can edit the config to leak environment variables |

Since version `<TODO release>`, Yor enforces two defense layers automatically:

### Layer 1 — env-var expansion is restricted to an allowlist

In `${env:NAME}` template values, only the following environment variable
names are expanded by default:

- Any name starting with `YOR_` (e.g. `YOR_OWNER`, `YOR_SIMPLE_TAGS`).
- The exact name `GIT_BRANCH` (preserved for backward compatibility with
  documented usage).

Every other name is left **literally unsubstituted** — Yor writes the raw
string `${env:NAME}` into the tag rather than the secret's value, and emits
a warning to the log. Well-known credential variables (`GITHUB_TOKEN`,
`GH_TOKEN`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`,
`AWS_SESSION_TOKEN`, `GOOGLE_APPLICATION_CREDENTIALS`,
`AZURE_CLIENT_SECRET`, `NPM_TOKEN`, `SLACK_TOKEN`) and any name containing
`SECRET`, `PASSWORD`, `PASSWD`, `PRIVATE_KEY`, `APIKEY`, `API_KEY`, or
`CREDENTIAL` are **always denied**, even if explicitly allowlisted.

Two environment variables tune this behavior:

- `YOR_ENV_ALLOWLIST` — comma-separated list of additional names (or
  `PREFIX*` globs) that operators can allow. Example:
  `YOR_ENV_ALLOWLIST=DEPLOY_REGION,MYORG_*`. The hard denylist above always
  wins, so this cannot be used to re-enable secret expansion.
- `YOR_DISABLE_ENV_EXPANSION=1` — global kill switch that disables
  `${env:...}` expansion entirely. Use this in hardened environments where
  no template expansion is needed.

### Layer 2 — the GitHub Action refuses to commit a repo-resident config

The official [`bridgecrewio/yor-action`](https://github.com/bridgecrewio/yor-action)
entrypoint refuses to auto-commit when `--config-file` resolves to a path
inside the directory passed to `-d`. In that case it emits a GitHub Actions
error annotation and exits with code `2`, so the workflow check fails and
nothing is pushed. To opt out, either:

- move the config file outside the scanned directory, or
- set the action input `commit_changes: 'NO'`.

This guard is purely a deployment-layer safeguard; running the Yor binary
directly (e.g. as a pre-commit hook or locally) is unaffected.
