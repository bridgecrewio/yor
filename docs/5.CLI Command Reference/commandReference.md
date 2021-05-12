# Cli Command Reference

The following parameter fields are used with the `./yor tag` command.

|Parameter  |  Description | Reference  |
|---|---|---|
|`-d DIRECTORY`, `--directory DIRECTORY` |IaC root directory. Can not be used together with --file. |[Scan Use Cases](doc:scan-use-cases#section-scan---repo-branch-folder-or-file)   |
|`--skip-tags` | run all tags except those listed | [Skipping Tags](https://github.com/bridgecrewio/yor/blob/yor-docs/docs/2.Using%20Yor/applyTag.md#skipping-tags) |
|`--skip-dirs` | run tags on all files within root directory except path listed | [Skipping Directories](https://github.com/bridgecrewio/yor/blob/yor-docs/docs/2.Using%20Yor/applyTag.md#skipping-directories)  |


The following parameter fields are used with the `./yor` command.

|Parameter  |  Description |
|---|---|
|`list-tags` |  list all the tags built into yor |
|`list-tag-groups` | list the groups of tags that are built into yor |

Type `yor -h` to have up-to-date list of supported commands.

