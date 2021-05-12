# Tagging Use Cases

## Tracing Code to Cloud

In order to locate your run time resource that was created based on a specific IaC resource - use ```yor_trace``` tag of this resource in build time and search this tag in relevant cloud console or CLI. Example below shows having such search from GitHub repository to AWS console.

![](../yor_trace.gif)

## Tracing Cloud to Code

In order to locate your IaC resource that created a given run time cloud resource, use ```git_file``` tag of this resource and navigate to relevant path. Alternatively - search ```yor_trace``` value in your repository. Example below shows having such search from AWS console to GitHub repository.

![](../yor_file.gif)

## Resource Ownership

Use ```git_last_modified_by``` and ```git_last_modified_At``` tags in your run time resources in order to determine who and when a resource last modified in associated IaC resource. Furthermore - discover list of all modifiers using ```git_modifiers``` tag. 

![](../yor_owner.gif)

