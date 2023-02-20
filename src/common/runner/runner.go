package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strconv"
	"strings"
	"sync"

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
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type Runner struct {
	TagGroups            []tagging.ITagGroup
	parsers              []common.IParser
	ChangeAccumulator    *reports.TagChangeAccumulator
	reportingService     *reports.ReportService
	dir                  string
	skipDirs             []string
	skippedTags          []string
	configFilePath       string
	skippedResourceTypes []string
	skippedResources     []string
	workersNum           int
	dryRun               bool
	localModuleTag       bool
}

const WorkersNumEnvKey = "YOR_WORKER_NUM"

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
	if commands.ConfigFile == "" {
		logger.Info("Did not get an external config file")
	}
	for _, tagGroup := range r.TagGroups {
		tagGroup.InitTagGroup(dir, commands.SkipTags, commands.Tag, tagging.WithTagPrefix(commands.TagPrefix))
		if simpleTagGroup, ok := tagGroup.(*simple.TagGroup); ok {
			simpleTagGroup.SetTags(extraTags)
		} else if externalTagGroup, ok := tagGroup.(*external.TagGroup); ok && commands.ConfigFile != "" {
			externalTagGroup.InitExternalTagGroups(commands.ConfigFile)
		}
	}
	processedParsers := map[string]struct{}{}
	for _, p := range commands.Parsers {
		if _, exists := processedParsers[p]; exists {
			continue
		}
		switch p {
		case "Terraform":
			r.parsers = append(r.parsers, &tfStructure.TerraformParser{})
		case "CloudFormation":
			r.parsers = append(r.parsers, &cfnStructure.CloudformationParser{})
		case "Serverless":
			r.parsers = append(r.parsers, &slsStructure.ServerlessParser{})
		default:
			logger.Warning(fmt.Sprintf("ignoring unknown parser %#v", err))
		}
		processedParsers[p] = struct{}{}
	}
	options := map[string]string{
		"tag-local-modules": strconv.FormatBool(commands.TagLocalModules)}
	for _, parser := range r.parsers {
		parser.Init(dir, options)
	}

	r.ChangeAccumulator = reports.TagChangeAccumulatorInstance
	r.reportingService = reports.ReportServiceInst
	r.dir = commands.Directory
	r.skippedTags = commands.SkipTags
	r.skipDirs = append(commands.SkipDirs, ".git")
	r.configFilePath = commands.ConfigFile
	r.dryRun = commands.DryRun
	if utils.InSlice(r.skipDirs, r.dir) {
		logger.Warning(fmt.Sprintf("Selected dir, %s, is skipped - expect an empty result", r.dir))
	}
	r.skippedResourceTypes = commands.SkipResourceTypes
	r.skippedResources = commands.SkipResources
	var convErr error
	r.workersNum, convErr = strconv.Atoi(utils.GetEnv(WorkersNumEnvKey, "10"))
	if convErr != nil {
		logger.Error(fmt.Sprintf("Got an invalid value for YOR_WORKERS_NUM, %v. If you didn't mean to leverage this option, please unset %v", os.Getenv(WorkersNumEnvKey), WorkersNumEnvKey))
	}
	return nil
}

func (r *Runner) toggleFinder(fileChan chan string, toggleExists chan bool, wg *sync.WaitGroup) {
	for file := range fileChan {
		if r.FindToggle(file) {
			toggleExists <- true
			r.AddSwitchVariable(file)
		}
		wg.Done()
	}
}

func (r *Runner) worker(fileChan chan string, wg *sync.WaitGroup, isToggleTurnedOn bool) {
	for file := range fileChan {
		r.TagFile(file, isToggleTurnedOn)
		wg.Done()
	}
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

	var finderWait sync.WaitGroup
	isToggleTurnedOn := false
	finderWait.Add(len(files))
	toggleExists := make(chan bool, r.workersNum)
	toggleFinderChan := make(chan string)

	for i := 0; i < r.workersNum; i++ {
		go r.toggleFinder(toggleFinderChan, toggleExists, &finderWait)
	}

	for _, file := range files {
		toggleFinderChan <- file
	}

	defer close(toggleExists)
	close(toggleFinderChan)
	finderWait.Wait()
	if len(toggleExists) != 0 {
		isToggleTurnedOn = <-toggleExists
	}

	var wg sync.WaitGroup
	wg.Add(len(files))
	fileChan := make(chan string)

	for i := 0; i < r.workersNum; i++ {
		go r.worker(fileChan, &wg, isToggleTurnedOn)
	}

	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)
	wg.Wait()

	for _, parser := range r.parsers {
		parser.Close()
	}

	return r.reportingService, nil
}

func (r *Runner) isSkippedResourceType(resourceType string) bool {
	for _, skippedResourceType := range r.skippedResourceTypes {
		if resourceType == skippedResourceType {
			return true
		}
	}
	return false
}

func (r *Runner) isSkippedResource(resource string) bool {
	for _, skippedResource := range r.skippedResources {
		if resource == skippedResource {
			return true
		}
	}
	return false
}

func (r *Runner) FindToggle(file string) bool {
	for _, parser := range r.parsers {
		if r.isFileSkipped(parser, file) {
			logger.Debug(fmt.Sprintf("%v parser Skipping %v", parser.Name(), file))
			continue
		}
		blocks, err := parser.ParseFile(file)
		if err != nil {
			logger.Info(fmt.Sprintf("Failed to parse file %v with parser %v", file, reflect.TypeOf(parser)))
			continue
		}
		for _, block := range blocks {
			terraformBlock, isTerraformType := block.(*tfStructure.TerraformBlock)
			if !isTerraformType {
				return false
			}

			if terraformBlock.HclSyntaxBlock.Type == tfStructure.VariableBlockType {
				if block.GetResourceID() == "yor_toggle" {
					rawBlock := block.GetRawBlock().(*hclwrite.Block)
					valueStr := string(rawBlock.Body().GetAttribute("default").Expr().BuildTokens(hclwrite.Tokens{}).Bytes())
					valueStr = strings.TrimSpace(valueStr)
					value, err := strconv.ParseBool(valueStr)
					if err != nil {
						logger.Info(fmt.Sprintf("Failed to convert default value of yor_toggle %v with parser %v", valueStr, err))
					}
					return value
				}
			}
		}
	}
	return false
}

func (r *Runner) AddSwitchVariable(file string) {
	src, err := os.ReadFile(file)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to read file %s because %s", file, err))
	}

	hclFile, diagnostics := hclwrite.ParseConfig(src, file, hcl.InitialPos)
	if diagnostics != nil && diagnostics.HasErrors() {
		hclErrors := diagnostics.Errs()
		logger.Info(fmt.Sprintf("Failed to parse hcl file %s because of errors %s", file, hclErrors))
	}
	rootBody := hclFile.Body()
	rootBody.AppendNewline()
	rootBody.AppendNewline()
	newBlock := rootBody.AppendNewBlock("variable", []string{"turn_off_yor_tags"})
	newBody := newBlock.Body()
	newBody.SetAttributeValue("default", cty.BoolVal(true))
	os.WriteFile(file, hclFile.Bytes(), 0644)
}

func (r *Runner) TagFile(file string, addToggle bool) {
	for _, parser := range r.parsers {
		if r.isFileSkipped(parser, file) {
			logger.Debug(fmt.Sprintf("%v parser Skipping %v", parser.Name(), file))
			continue
		}
		logger.Info(fmt.Sprintf("Tagging %v\n", file))
		blocks, err := parser.ParseFile(file)
		if err != nil {
			logger.Info(fmt.Sprintf("Failed to parse file %v with parser %v", file, reflect.TypeOf(parser)))
			continue
		}
		isFileTaggable := false
		for _, block := range blocks {
			if r.isSkippedResourceType(block.GetResourceType()) {
				continue
			}
			if r.isSkippedResource(block.GetResourceID()) {
				continue
			}
			if block.IsBlockTaggable() {
				logger.Debug(fmt.Sprintf("Tagging %v:%v", file, block.GetResourceID()))
				isFileTaggable = true
				for _, tagGroup := range r.TagGroups {
					err := tagGroup.CreateTagsForBlock(block)
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to tag %v in %v due to %v", block.GetResourceID(), block.GetFilePath(), err.Error()))
						continue
					}
				}
			} else {
				logger.Debug(fmt.Sprintf("Block %v:%v is not taggable, skipping", file, block.GetResourceID()))
			}
			r.ChangeAccumulator.AccumulateChanges(block)
		}
		if isFileTaggable && !r.dryRun {
			err = parser.WriteFile(file, blocks, file, addToggle)
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
	relPath, _ := filepath.Rel(r.dir, file)
	for _, sp := range r.skipDirs {
		if strings.HasPrefix(r.dir+"/"+relPath, sp) {
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
	return !p.ValidFile(file)
}
