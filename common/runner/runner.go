package runner

import (
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/common/logger"
	"bridgecrewio/yor/common/reports"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
	"bridgecrewio/yor/common/tagging/tags"
	tfStructure "bridgecrewio/yor/terraform/structure"
	tfTagging "bridgecrewio/yor/terraform/tagging"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

type Runner struct {
	taggers           []tagging.ITagger
	parsers           []structure.IParser
	gitService        *gitservice.GitService
	changeAccumulator *reports.TagChangeAccumulator
	reportingService  *reports.ReportService
}

func (r *Runner) Init(dir string, extraTagsFromArgs map[string]string, externalTagsPath string) error {
	gitService, err := gitservice.NewGitService(dir)
	if err != nil {
		logger.Error("Failed to initialize git service")
	}
	r.gitService = gitService
	r.taggers = append(r.taggers, &tfTagging.TerraformTagger{})
	extraTags, err := loadExternalTags(externalTagsPath)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to load extenal tags from plugins due to error: %s", err))
	}
	extraTags = append(extraTags, createCmdTags(extraTagsFromArgs)...)
	for _, tagger := range r.taggers {
		tagger.InitTags(extraTags)
	}
	r.parsers = append(r.parsers, &tfStructure.TerrraformParser{})
	for _, parser := range r.parsers {
		parser.Init(dir, nil)
	}

	r.changeAccumulator = reports.TagChangeAccumulatorInstance
	r.reportingService = reports.ReportServiceInst
	return nil
}

func (r *Runner) TagDirectory(dir string) (*reports.Report, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("Failed to scan dir", path)
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		logger.Error("Failed to run Walk() on root dir", dir)
	}

	for _, file := range files {
		r.TagFile(file)
	}
	//	TODO - return Report's result from this method
	reportService := reports.ReportService{}

	return reportService.CreateReport(), nil
}

func (r *Runner) TagFile(file string) {
	for _, parser := range r.parsers {
		blocks, err := parser.ParseFile(file)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to parse file %v with parser %v", file, parser))
			continue
		}
		for _, block := range blocks {
			for _, tagger := range r.taggers {
				if block.IsBlockTaggable() {
					blame, err := r.gitService.GetBlameForFileLines(file, block.GetLines())
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to tag %v with git tags, err: %v", block.GetResourceID(), err.Error()))
						continue
					}
					tagger.CreateTagsForBlock(block, blame)
					r.changeAccumulator.AccumulateChanges(block)
				}
			}
			err = parser.WriteFile(file, blocks)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed writing tags to file %s, because %v", file, err))
			}
			//	TODO: if block is a local module, run TagDir on it as well
			//  Need to avoid cycles here!!
		}
	}
	r.reportingService.CreateReport()

	// TODO: support multiple output formats according to args
	r.reportingService.PrintToStdout()
}

func createCmdTags(extraTagsFromArgs map[string]string) []tags.ITag {
	extraTags := make([]tags.ITag, len(extraTagsFromArgs))
	index := 0
	for key := range extraTagsFromArgs {
		newTag := tags.Init(key, extraTagsFromArgs[key])
		extraTags[index] = newTag
		index++
	}

	return extraTags
}

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
