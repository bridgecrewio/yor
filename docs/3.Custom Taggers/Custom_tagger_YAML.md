---
layout: default
published: true
title: Custom Tagger Using YAML  Configuration Files
nav_order: 2
---
# Custom Tagger Using YAML  Configuration Files

The Yor framework uses YAML configuration files to support advanced rules when applying custom tags.
Users can define tagging enforcement rules that are specific to their organization’s needs. 
YAML based custom tagging enables you to have different tags for different resource types and existing resource tags.

## Running YAML based
In the CLI, define the path of the YAML configuration file that you want to apply. For example:

`yor tag --config-file </path/to/file>`

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
    1. *Resource_types* sequence: resource types to tag (Terraform format)
    2. *Tags sequence*: tag resources that have all the tags within the map
4. Directory definition (optional): path to defined taggable resources

***Example 2:** Tagging specific resource types with specific `key:value` tags in a defined directory.*
```
name: env
value:
    default: prod
filters:
    resource_types:
        - aws_ec2_instance
        - aws_security_group
    tags:
        git_modifiers: donnaj
        git_repo: checkov
    directory: /path/to/some/dir
```

5. Use case dynamic value definition using *value* mapping (optional): Tags are defined based on matching
   keys that contain a sequence of values. Under each value the user can define which existing tags a resource will be 
   tagged with. If none of the conditions are matched, a default value will be applied. In the example below
   `- aws_ec2_instance` and `aws_security_group` in the directory `/path/to/some/dir` and existing tag `yor_trace: 123` will 
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
   filters:
        resource_types:
              - aws_ec2_instance
              - aws_security_group
   tags:
        yor_trace: 123
   directory: /path/to/some/dir
```
## Custom tagging using CLI

You can use some YAML configuration capabilities in a CLI command. 
1. `--tag-name`: define tag name
2. `--tag-value`: define tag value
3. `--resource-types`: define which resource types to tag (Terraform format)
4. `-filter-tags`: tag resources that have tags as defined. Use an array [] to support multiple values and to support `AND` logic between tags

In the example below, EC2 instances and Security Groups will be tagged with the `env:prod` tag. Use this in cases where a resource that has `tronxd` 
or `amy` are one of the `git_modifiers` and it is located in `checkov` or `terragoat git_repo`.

**Example 3:** CLI custom tagging

```yor tag --tag-name env –tag-value prod --resource-types [aws_ec2_instance,aws_ec2_security_group] –filter-tags git_modifiers=[tronxd,amy];git_repo=[checkov,terragoat]```

## Running Yor with Custom Taggers
Use the following example to run Yor with custom tags:
```sh
./yor tag --custom-tagging tests/yor_plugins/example
# run yor with custom tags located in tests/yor_plugins/example

./yor tag --custom-tagging tests/yor_plugins/example,tests/yor_plugins/tag_group_example
# run yor with custom tags located in tests/yor_plugins/example and custom taggers located in tests/yor_plugins/tag_group_example
```
