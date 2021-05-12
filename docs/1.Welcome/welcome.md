# Overview
Yor is an open-source tool that helps to manage tags in a consistent manner across infrastructure as code frameworks 
(Terraform, Cloudformation, Kubernetes, and Serverless Framework). By auto-tagging in IaC you will be able to trace any cloud resource from code to cloud.

# IaC Types
Yor enables version-controlled owner assignment and resource tracing based git history. It also can extend tag enforcement logic by loading external tagging logic into the CI/CD pipeline. 

Yor can tag these IaC file types:
  * Terraform (for AWS, GCP and Azure)
  * CloudFormation (YAML, JSON)
  * Serverless
  * K8S (YAML, JSON)

# Tracing Tagger
* ```yor_trace``` tag enables full attribution between build time and run time resources. 

# Git-based Tagger
* ```git_*``` tags connect cloud resources to individual git commits and enable assigning clear ownership between developers and the resources they routinely change.

Yor collects data from [git-blame](https://git-scm.com/docs/git-blame) and enables mapping individual resources to specific commits.

### Supported tags

```
git_org = "bridgecrewio"
git_repo = "terragoat"
git_file = "README.md" # this is the path from the repo root dir...
git_commit = "47accf06f13b503f3bab06fed7860e72f7523cac" # This is the latest commit for this resource
git_last_modified_At = "2020-03-28 21:42:46 +0000 UTC"
git_last_modified_by = "schosterbarak@gmail.com"
git_modifiers = "schosterbarak/baraks" # These are extracted from the emails, everything before the @ sign. Can be done for modified_by tag as well
```
