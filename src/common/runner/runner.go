package runner

import (
	cfnStructure "bridgecrewio/yor/src/cloudformation/structure"
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/reports"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/code2cloud"
	"bridgecrewio/yor/src/common/tagging/gittag"
	"bridgecrewio/yor/src/common/tagging/simple"
	"bridgecrewio/yor/src/common/tagging/tags"
	tfStructure "bridgecrewio/yor/src/terraform/structure"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

type Runner struct {
	taggers           []tagging.ITagger
	parsers           []structure.IParser
	changeAccumulator *reports.TagChangeAccumulator
	reportingService  *reports.ReportService
	dir               string
	skipDirs          []string
	skippedTags       []string
}

func (r *Runner) Init(commands *common.Options) error {
	dir := commands.Directory
	r.taggers = append(r.taggers, &gittag.Tagger{}, &simple.Tagger{}, &code2cloud.Tagger{})
	for _, tagger := range r.taggers {
		tagger.InitTagger(dir, commands.SkipTags)
		if simpleTagger, ok := tagger.(*simple.Tagger); ok {
			extraTags, err := loadExternalResources(commands.CustomTagging)
			if err != nil {
				logger.Warning(fmt.Sprintf("failed to load extenal tags from plugins due to error: %s", err))
			} else {
				extraTags = append(extraTags, createCmdTags(commands.ExtraTags)...)
				simpleTagger.SetTags(extraTags)
			}
		}
	}
	r.parsers = append(r.parsers, &tfStructure.TerrraformParser{}, &cfnStructure.CloudformationParser{})
	for _, parser := range r.parsers {
		parser.Init(dir, nil)
	}

	r.changeAccumulator = reports.TagChangeAccumulatorInstance
	r.reportingService = reports.ReportServiceInst
	r.dir = commands.Directory
	r.skippedTags = commands.SkipTags
	r.skipDirs = commands.SkipDirs

	if common.InSlice(r.skipDirs, r.dir) {
		logger.Warning(fmt.Sprintf("Selected dir, %s, is skipped - expect an empty result", r.dir))
	}
	return nil
}

func (r *Runner) TagDirectory() (*reports.ReportService, error) {
	var files []string
	err := filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("Failed to scan dir", path)
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		logger.Error("Failed to run Walk() on root dir", r.dir)
	}

	for _, file := range files {
		r.TagFile(file)
	}

	return r.reportingService, nil
}

func (r *Runner) TagFile(file string) {
	for _, parser := range r.parsers {
		if r.isFileSkipped(parser, file) {
			logger.Info("Skipping", file)
			continue
		}
		blocks, err := parser.ParseFile(file)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to parse file %v with parser %v", file, parser))
			continue
		}
		isFileTaggable := false
		for _, block := range blocks {
			if block.IsBlockTaggable() {
				isFileTaggable = true
				for _, tagger := range r.taggers {
					tagger.CreateTagsForBlock(block)
				}
			}
			r.changeAccumulator.AccumulateChanges(block)
		}
		if isFileTaggable {
			err = parser.WriteFile(file, blocks, file)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed writing tags to file %s, because %v", file, err))
			}
		}
	}

}

func createCmdTags(extraTagsStr string) []tags.ITag {
	var extraTagsFromArgs map[string]string
	if err := json.Unmarshal([]byte(extraTagsStr), &extraTagsFromArgs); err != nil {
		logger.Error(fmt.Sprintf("failed to parse extra tags: %s", err))
	}
	extraTags := make([]tags.ITag, len(extraTagsFromArgs))
	index := 0
	for key := range extraTagsFromArgs {
		newTag := tags.Init(key, extraTagsFromArgs[key])
		extraTags[index] = newTag
		index++
	}

	return extraTags
}

func loadExternalResources(externalPaths []string) ([]tags.ITag, error) {
	var extraTags []tags.ITag
	var extraTaggers []tagging.ITagger
	var plugins []string

	for _, path := range externalPaths {
		// find all .so files under the given externalPaths
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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

			iPtrs, err := extractExternalResources(plug, "ExtraTags")
			if err != nil {
				return nil, err
			}
			for _, iTag := range iPtrs {
				tag, ok := iTag.(tags.ITag)
				if !ok {
					return nil, fmt.Errorf("unexpected type from module symbol")
				}
				extraTags = append(extraTags, tag)
			}
			iPtrs, err = extractExternalResources(plug, "ExtraTaggers")
			for _, iTagger := range iPtrs {
				if tagger, ok := iTagger.(tagging.ITagger); ok {
					extraTaggers = append(extraTaggers, tagger)
				} else {
					return nil, fmt.Errorf("unexpected type from module symbol ExtraTaggers")
				}
			}
		}
	}

	return extraTags, nil
}

func extractExternalResources(plug *plugin.Plugin, symbol string) ([]interface{}, error) {
	symExtraTags, err := plug.Lookup(symbol)
	if err != nil {
		return nil, nil
	}
	logger.Info("Found values for the symbol:", symbol)
	// convert to its actual type, *[]interface{}
	var iArrPtr *[]interface{}
	iArrPtr, ok := symExtraTags.(*[]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type from module symbol")
	}

	return *iArrPtr, nil
}

func (r *Runner) isFileSkipped(p structure.IParser, file string) bool {
	for _, sp := range r.skipDirs {
		if strings.HasPrefix(file, sp) {
			return true
		}
	}

	matchingSuffix := false
	for _, suffix := range p.GetAllowedFileTypes() {
		if strings.HasSuffix(file, suffix) {
			matchingSuffix = true
		}
	}
	if !matchingSuffix {
		return true
	}
	for _, pattern := range p.GetSkippedDirs() {
		if strings.Contains(file, pattern) {
			return true
		}
	}
	return false
}
