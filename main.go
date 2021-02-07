package main

import (
	"bridgecrewio/yor/common/reports"
	"bridgecrewio/yor/common/tagging/tags"
	"fmt"
)

func main() {
	fmt.Println("Welcome to Yor!")
}

func parseArgs(args ...interface{}) {
	// TODO
}

func printReport(report *reports.Report) {
	// TODO
}

func createExtraTags(extraTagsFromArgs map[string]string) []tags.ITag {
	extraTags := make([]tags.ITag, len(extraTagsFromArgs))
	index := 0
	for key := range extraTagsFromArgs {
		newTag := tags.Init(key, extraTagsFromArgs[key])
		extraTags[index] = newTag
		index++
	}

	return extraTags
}
