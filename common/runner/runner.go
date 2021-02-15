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
)

type Runner struct {
	taggers           []tagging.ITagger
	parsers           []structure.IParser
	gitService        *gitservice.GitService
	changeAccumulator *reports.TagChangeAccumulator
	reportingService  *reports.ReportService
}

func (r *Runner) Init(dir string, extraTags []tags.ITag) error {
	gitService, err := gitservice.NewGitService(dir)
	if err != nil {
		logger.Error("Failed to initialize git service")
	}
	r.gitService = gitService
	r.taggers = append(r.taggers, &tfTagging.TerraformTagger{})
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
			parser.WriteFile(file, blocks)
			//	TODO: if block is a local module, run TagDir on it as well
			//  Need to avoid cycles here!!
		}
	}
	r.reportingService.CreateReport()

	// TODO: support multiple output formats according to args
	r.reportingService.PrintToStdout()
}
