---
layout: default
published: true
title: Reviewing Tag Results
nav_order: 4
---

# Reviewing Tag Results
You can assume that each resource that is being tagged using Yor has a diff view. For example -  

![Diff View](/yor_diff_view.png)

After applying `./yor tag` command, you will get the CLI Findings Summary. This is also available once you are running 
`./yor tag -d . -o cli` 

![Yor Summary](/yor_summary.png)

Use the following commands to display the tags that are currently used.
```sh
./yor tag -d . -o json
# json output

./yor tag -d . --output cli --output-json-file result.json
# will print cli output and additional output to file on json file -- enables programmatic analysis alongside printing human readable result
```

For a JSON file example see 

![YOR JSON Results](/yor_json_results.png)




