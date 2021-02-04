package main

import (
	"bridgecrewio/yor/common/reports"
	"bridgecrewio/yor/common/tagging"
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

func createExtraTags(extraTagsFromArgs map[string]interface{}) []tagging.ITag {
	extraTags := make([]tagging.ITag, len(extraTagsFromArgs))
	index := 0
	for key := range extraTagsFromArgs {
		newTag := tagging.Init(key, extraTagsFromArgs[key])
		extraTags[index] = newTag
		index++
	}

	return extraTags
}
