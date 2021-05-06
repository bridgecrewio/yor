# Supported Taggers

The following commands are used to see the list of Yor supported taggers. 
`list-tag-groups` - list the groups of tags that are built into yor
   ```sh
   ./yor list-tag-groups
   ```
`list-tags` - lists all the tags. This will print each tag key
   and the relevant group the tag belongs to.
   ```sh
    ./yor list-tags 
    # List all the tags built into yor
   ./yor list-tags --tag-groups git
    # List all the tags built into yor under the tag group git
    ```
