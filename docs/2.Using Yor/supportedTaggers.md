---
layout: default
published: true
title: Supported Taggers
nav_order: 2
---
# Viewing Supported Taggers in Your Environment

The following commands are used to see the list of Yor supported taggers, both [built-in](../1.Welcome/welcome.md#built-in-taggers) and [custom tags](../3.Custom%20Taggers/customTagExamples.md). 

`list-tag-groups` - list the groups of tags that are built into yor
   ```sh
   ./yor list-tag-groups
   ```
`list-tags` - lists all the tags. This will print each tag key, and the relevant group the tag belongs to.
   ```sh
    ./yor list-tags 
    # List all the tags built into yor
   ```
![Environment variables after tagging](../yor_list_tags_after_env_var.png)
   
   
   
   ```sh
   ./yor list-tags --tag-groups git
    # List all the tags built into yor under the tag group git
    ```
