# Overwriting Source Files Algorithm

To avoid overwriting existing values, these are the steps the writer should take:
1. Copy over whatever the file contains before the resource lines
2. For every resource

    a. If the resource does not have tags - create new tags

    b. If the resource has tags: (1) update the existing modified tags in place and (2) add the new tags at the end

![Diagram](https://lucid.app/publicSegments/view/7de70ec0-e9ec-4dcf-8612-f2bf22f8e1bc/image.png)
