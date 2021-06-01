package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"

	cfnStructure "github.com/bridgecrewio/yor/src/cloudformation/structure"
	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/clioptions"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/reports"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/external"
	"github.com/bridgecrewio/yor/src/common/tagging/simple"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	taggingUtils "github.com/bridgecrewio/yor/src/common/tagging/utils"
	"github.com/bridgecrewio/yor/src/common/utils"
	slsStructure "github.com/bridgecrewio/yor/src/serverless/structure"
	tfStructure "github.com/bridgecrewio/yor/src/terraform/structure"
)

type Runner struct {
	TagGroups         []tagging.ITagGroup
	parsers           []common.IParser
	ChangeAccumulator *reports.TagChangeAccumulator
	reportingService  *reports.ReportService
	dir               string
	skipDirs          []string
	skippedTags       []string
	configFilePath    string
}

func (r *Runner) Init(commands *clioptions.TagOptions) error {
	dir := commands.Directory
	extraTags, extraTagGroups, err := loadExternalResources(commands.CustomTagging)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to load extenal tags from plugins due to error: %s", err))
	}
	for _, group := range commands.TagGroups {
		tagGroup := taggingUtils.TagGroupsByName(taggingUtils.TagGroupName(group))
		r.TagGroups = append(r.TagGroups, tagGroup)
	}
	r.TagGroups = append(r.TagGroups, extraTagGroups...)
	for _, tagGroup := range r.TagGroups {
		tagGroup.InitTagGroup(dir, commands.SkipTags)
		if simpleTagGroup, ok := tagGroup.(*simple.TagGroup); ok {
			simpleTagGroup.SetTags(extraTags)
		} else if externalTagGroup, ok := tagGroup.(*external.TagGroup); ok && commands.ConfigFile != "" {
			externalTagGroup.InitExternalTagGroups(commands.ConfigFile)
		}
	}
	r.parsers = append(r.parsers, &tfStructure.TerrraformParser{}, &cfnStructure.CloudformationParser{}, &slsStructure.ServerlessParser{})
	for _, parser := range r.parsers {
		parser.Init(dir, nil)
	}

	r.ChangeAccumulator = reports.TagChangeAccumulatorInstance
	r.reportingService = reports.ReportServiceInst
	r.dir = commands.Directory
	r.skippedTags = commands.SkipTags
	r.skipDirs = commands.SkipDirs
	r.configFilePath = commands.ConfigFile
	if utils.InSlice(r.skipDirs, r.dir) {
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
		logger.Info("Tagging %v\n", file)
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
			logger.Info(fmt.Sprintf("Failed to parse file %v with parser %v", file, reflect.TypeOf(parser)))
			continue
		}
		isFileTaggable := false
		for _, block := range blocks {
			if block.IsBlockTaggable() {
				isFileTaggable = true
				for _, tagGroup := range r.TagGroups {
					err := tagGroup.CreateTagsForBlock(block)
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to tag %v in %v due to %v", block.GetResourceID(), block.GetFilePath(), err.Error()))
						continue
					}
				}
			}
			r.ChangeAccumulator.AccumulateChanges(block)
		}
		if isFileTaggable {
			err = parser.WriteFile(file, blocks, file)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed writing tags to file %s, because %v", file, err))
			}
		}
	}

}

func loadExternalResources(externalPaths []string) ([]tags.ITag, []tagging.ITagGroup, error) {
	var extraTags []tags.ITag
	var extraTagGroups []tagging.ITagGroup
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
			return nil, nil, err
		}

		for _, pluginPath := range plugins {
			plug, err := plugin.Open(pluginPath)
			if err != nil {
				return nil, nil, err
			}

			iPtrs, err := extractExternalResources(plug, "ExtraTags")
			if err != nil {
				return nil, nil, err
			}
			for _, iTag := range iPtrs {
				tag, ok := iTag.(tags.ITag)
				if !ok {
					return nil, nil, fmt.Errorf("unexpected type from module symbol")
				}
				extraTags = append(extraTags, tag)
			}
			iPtrs, err = extractExternalResources(plug, "ExtraTagGroups")
			if err != nil {
				return nil, nil, err
			}
			for _, iTagGroup := range iPtrs {
				if tagGroup, ok := iTagGroup.(tagging.ITagGroup); ok {
					extraTagGroups = append(extraTagGroups, tagGroup)
				} else {
					return nil, nil, fmt.Errorf("unexpected type from module symbol ExtraTagGroups")
				}
			}
		}
	}

	return extraTags, extraTagGroups, nil
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

func (r *Runner) isFileSkipped(p common.IParser, file string) bool {
	for _, sp := range r.skipDirs {
		if strings.HasPrefix(file, sp) {
			return true
		}
	}

	matchingSuffix := false
	for _, suffix := range p.GetSupportedFileExtensions() {
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
