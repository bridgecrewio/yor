package main

import (
	"bridgecrewio/yor/common/logger"
	"bridgecrewio/yor/common/tagging/tags"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

func main() {
	fmt.Println("Welcome to Yor!")
}

// TODO
// func parseArgs(args ...interface{}) {
// }

// TODO
// func printReport(report *reports.Report) {
// }

// func createExtraTags(extraTagsFromArgs map[string]string) []tags.ITag {
//	extraTags := make([]tags.ITag, len(extraTagsFromArgs))
//	index := 0
//	for key := range extraTagsFromArgs {
//		newTag := tags.Init(key, extraTagsFromArgs[key])
//		extraTags[index] = newTag
//		index++
//	}
//
//	return extraTags
//}

func loadExternalTags(tagsPath string) ([]tags.ITag, error) {
	var extraTags []tags.ITag
	var plugins []string

	// find all .so files under the given tagsPath
	err := filepath.Walk(tagsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".so") {
			plugins = append(plugins, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, pluginPath := range plugins {
		plug, err := plugin.Open(pluginPath)
		if err != nil {
			return nil, err
		}

		// extract the symbol "ExtraTags" from the plugin file
		symExtraTags, err := plug.Lookup("ExtraTags")
		if err != nil {
			logger.Warning(err.Error())
			continue
		}

		// convert ExtraTags to its actual type, *[]interface{}
		var iTagsPtr *[]interface{}
		iTagsPtr, ok := symExtraTags.(*[]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected type from module symbol")
		}

		iTags := *iTagsPtr
		for _, iTag := range iTags {
			tag, ok := iTag.(tags.ITag)
			if !ok {
				return nil, fmt.Errorf("unexpected type from module symbol")
			}
			extraTags = append(extraTags, tag)
		}
	}

	return extraTags, nil
}
