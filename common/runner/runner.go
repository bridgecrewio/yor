package runner

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/reports"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
	"bridgecrewio/yor/common/tagging/tags"
	tfStructure "bridgecrewio/yor/terraform/structure"
	tfTagging "bridgecrewio/yor/terraform/tagging"
	"os"
	"path/filepath"
)

type Runner struct {
	taggers    []tagging.ITagger
	parsers    []structure.IParser
	gitService *git_service.GitService
}

func NewRunner(taggerTypes []tagging.TaggerType, extraTags []tags.ITag) *Runner {
	// TODO
	return nil
}

func initTaggers(taggerTypes []tagging.TaggerType, extraTags []tags.ITag) {
	// TODO: create a new Tagger instance and send the extraTags as param
}

func (r *Runner) Init(dir string) error {
	gitService, err := git_service.NewGitService(dir)
	if err != nil {
		return err
	}
	r.gitService = gitService
	r.taggers = append(r.taggers, &tfTagging.TerraformTagger{})
	r.parsers = append(r.parsers, &tfStructure.TerrraformParser{})
	for _, parser := range r.parsers {
		parser.Init()
	}

	return nil
}

func (r *Runner) TagDirectory(dir string) (*reports.Report, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		err = r.TagFile(file, dir)
		if err != nil {
			return nil, err
		}
	}
	//	TODO - return Report's result from this method
	reportService := reports.ReportService{}

	return reportService.CreateReport(), nil
}

func (r *Runner) TagFile(file string, rootDir string) error {
	for _, parser := range r.parsers {
		blocks, err := parser.ParseFile(file, rootDir)
		if err != nil {
			return err
		}
		for _, block := range blocks {
			for _, tagger := range r.taggers {
				if block.IsBlockTaggable() {
					blame, err := r.gitService.GetBlameForFileLines(file, block.GetLines())
					if err != nil {
						return err
					}
					tagger.CreateTagsForBlock(block, blame)
				}
			}
			//	TODO: if block is a local module, run TagDir on it as well
			//  Need to avoid cycles here!!
		}
	}
	return nil
}
