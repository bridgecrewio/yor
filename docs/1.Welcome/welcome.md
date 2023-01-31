---
layout: default
published: true
title: What is Yor?
nav_order: 1
---

# Overview

Yor is an open-source tool that helps to manage tags consistently across infrastructure as code (IaC) frameworks. Auto-tagging in IaC enables you to trace any resource from code to cloud.

IaC resource tagging provides added benefits which include:
* **Security Risk Management** - identify resources that require heightened security
* **Access Control** - use IAM policies to control access to resources
* **Automation** - tag resources for automation policies
* **Cost Allocation** - break down costs based on tagging of resources
* **Operation Support** - tracing of resource ownership
* **Console Organization** - consolidate resource inventory based on tags so that you can track multiple instances of the same resources

# IaC Types
Yor enables version-controlled owner assignment and resource tracing based on git history and extends tag enforcement logic
into the CI/CD pipeline.

Yor can tag the following IaC file types:
  * Terraform (for AWS, GCP, Azure, AliCloud, OCI)
  * CloudFormation (YAML, JSON)
  * Serverless Framework

# Built-in Taggers

Yor supports a number of different taggers that can be used to trace and identify resources. The following types of tags are
included with the Yor distribution.

## Tracing Tagger

Yor provides a unique ID created when running the tag command which enables complete traceability between build time and run time resources.
The ability to track a runtime resource relies on a resource block in your IaC files creating better visibility of your assets. Yor can also detect
drifts between build time and runtime resources across IaC frameworks and
multiple cloud providers.

### Supported Yor Trace Tags
The following tags are supported in Yor:

```yor_trace``` is a unique ID provided when a resource is tagged.

For examples see [Use Cases](../4.Use Cases/useCases.md).

## Git-based Tagger
Yor collects data from [git-blame](https://git-scm.com/docs/git-blame) logs to create tags which enable the mapping of individual
resources to specific commits.

The ```git_*``` tags connect cloud resources to individual git commits and establish clear ownership between developers and
resources which are routinely change.

### Supported Git tags
The following tags are supported in Yor:
```
git_org = "bridgecrewio"
git_repo = "terragoat"
git_file = "README.md" # This is the path from the repo root dir.
git_commit = "47accf06f13b503f3bab06fed7860e72f7523cac" # This is the latest commit for this resource.
git_last_modified_At = "2020-03-28 21:42:46 +0000 UTC"
git_last_modified_by = "schosterbarak@gmail.com"
git_modifiers = "schosterbarak/baraks" # These are extracted from emails (everything before the @ sign). This can also be done
for the git_last_modified_by tag.
```

# Custom Taggers

Yor supports Custom taggers and tag groups to enable you to enhance your resource traceability. Yor supports custom taggers using:
* [Environment variable settings](../3.Custom Taggers/customTagExamples.md#adding-simple-tags-using-environment-variables)
* [Golang settings](../3.Custom Taggers/customTagExamples.md#adding-custom-tags-using-golang)
* [YAML configuration files](../3.Custom Taggers/Custom_tagger_YAML.md#custom-tagger-using-yaml-configuration-files)
* [CLI commands](../3.Custom Taggers/Custom_tagger_YAML.md#custom-tagging-using-cli)

Using custom tags provides organizations with the ability to tag resources to match the development cycle, development flow, or the organization's
structure.
