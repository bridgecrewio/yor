---
layout: default
published: true
title: CLI Command Reference
nav_order: 1
---

# Cli Command Reference

The following parameter fields are used with the `./yor tag` command.

|Parameter  |  Description | Reference  |
|---|---|---|
|`-d DIRECTORY`, `--directory DIRECTORY` |IaC root directory. Can not be used together with --file. |   |
|`--skip-tags` | run all tags except those listed | [Skipping Tags](../2.Using Yor/applyTag.md#skipping-tags) |
|`--skip-dirs` | run tags on all files within root directory except path listed | [Skipping Directories](../2.Using Yor/applyTag.md#skipping-directories)  |


The following parameter fields are used with the `./yor` command.

|Parameter  |  Description |
|---|---|
|`list-tags` |  list all the tags built into yor |
|`list-tag-groups` | list the groups of tags that are built into yor |

Type `yor -h` to have up-to-date list of supported commands.
