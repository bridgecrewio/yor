---
layout: default
published: true
title: Reviewing Tag Results
nav_order: 4
---

# Reviewing Tag Results
You can assume that each resource that is being tagged using Yor has a diff view such of the following:
<Yor_diff_view.png>

After applying `./yor tag` command - user will get in CLI Findings Summary, which is also available once running `./yor tag -d . -o cli`
<Yor_summary.png>

Use the following commands to display the tags that are currently used.
```sh
./yor tag -d . -o json
# json output

./yor tag -d . --output cli --output-json-file result.json
# will print cli output and additional output to file on json file -- enables programatic analysis alongside printing human readable result
```

json file example:
<Yor_json_results.png>




