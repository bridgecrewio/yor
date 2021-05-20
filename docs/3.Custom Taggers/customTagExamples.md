# Examples of Applied Custom Tags

## Prerequisites

Examples can be found in `tests/yor_plugins`
Yor supports 2 ways of adding custom tags:
1. [Using Enverionment Variable - Simple tags with constant key-value](#adding-simple-tags)
2. Using Go builtim classes -
	1. [Simple code-based tags](#adding-simple-code-based-tags)
	2. [Complex tags which rely on different inputs](#adding-complex-tags)

## Adding Simple Tags Using Enviornment Variables
To add tags with constant key-value pairs, set the environment variable YOR_SIMPLE_TAGS
with a JSON object detailing all key value pairs. Example:
```sh
export YOR_SIMPLE_TAGS='{"somekey": "somevalue", "another-key": "another_val"}'
# When yor is run, all resources will be tagged by these two tags as well
```

## Adding Custom Tags Using Golang

### Adding Simple Code Based Tags Using Golang
1. Create tags implementing the `ITag` interface.
2. If you wish to override an existing tag, make the tag's method `GetPriority()` return a positive number. Otherwise, return `0` or a negative number.
3. Create a file located in package `main` that exposes a variable `ExtraTags` - array containing pointers to all tags implemented. Example:
    ```go
    package main
    
    var ExtraTags = []interface{}{&TerragoatTag{}, &CheckovTag{}}
    ```
4. Run the command `go build -gcflags="all=-N -l" -buildmode=plugin -o <plugin-dir>/extra_tags.so <plugin-dir>/*.go`

See example in [tests/yor_plugins/example](tests/yor_plugins/example)

### Adding Complex Tags Using Golang
1. Create a tagger struct, implementing the `ITagGroup` interface.
2. Implement the `InitTagGroup` method, which should look something like this:
    ```go
    func (d *CustomTagger) InitTagGroup(_ string, skippedTags []string) {
	    d.SkippedTags = skippedTags
	    d.SetTags([]tags.ITag{}) // This is just a placeholder
    }
    ```
3. Implement the `CreateTagsForBlock` method, which will look something like this:
    ```go
   func (d *CustomTagger) CreateTagsForBlock(block structure.IBlock) {
        var newTags []tags.ITag
        for _, tag := range d.GetTags() {
            tagVal, err := tag.CalculateValue(<Whichever struct you choose to pass to the tagger>)
            if err != nil {
                logger.Error(fmt.Sprintf("Failed to create %v tag for block %v", tag.GetKey(), block.GetResourceID()))
            }
            newTags = append(newTags, tagVal)
        }
        block.AddNewTags(newTags)
   }
    ```
4. Implement tags, which implement the `ITag` interface just like we described [here](#adding-simple-tags)
5. Go back to the `InitTagGroup` method and add pointers to your new tags in the input of the `SetTags` function call.
6. Create a file located in package `main` that exposes a variable `ExtraTaggers` - array containing pointers to all tags implemented. Example:
    ```go
    package main
    
    var ExtraTaggers = []interface{}{&CustomTagger{}}
    ```

See example in [tests/yor_plugins/tag_group_example](tests/yor_plugins/tag_group_example)

## Running yor with the external tags / taggers

```sh
./yor tag --custom-tagging tests/yor_plugins/example
# run yor with custom tags located in tests/yor_plugins/example

./yor tag --custom-tagging tests/yor_plugins/example,tests/yor_plugins/tag_group_example
# run yor with custom tags located in tests/yor_plugins/example and custom taggers located in tests/yor_plugins/tag_group_example
```

