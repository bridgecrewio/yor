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
