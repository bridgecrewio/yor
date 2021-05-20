---
layout: default
published: true
title: Reviewing Tag Results
nav_order: 4
---

# Reviewing Tag Results
Use the following commands to display the tags that are currently used.
```sh
./yor tag -d . -o cli
# default cli output

./yor tag -d . -o json
# json output

./yor tag -d . --output cli --output-json-file result.json
# will print cli output and additional output to file on json file -- enables programatic analysis alongside printing human readable result
```