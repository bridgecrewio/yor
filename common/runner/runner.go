package runner

import (
	"bridgecrewio/yor/common/reports"
	"bridgecrewio/yor/common/tagging"
	"bridgecrewio/yor/common/tagging/tags"
)

type Runner struct {
	taggers []tagging.ITagger
}

func NewRunner(taggerTypes []tagging.TaggerType, extraTags []tags.ITag) *Runner {
	// TODO
	return nil
}

func initTaggers(taggerTypes []tagging.TaggerType, extraTags []tags.ITag) {
	// TODO: create a new Tagger instance and send the extraTags as param
}

func (r *Runner) TagDirectory(dir string) *reports.Report {
	//	TODO - return Report's result from this method
	reportService := reports.ReportService{}

	return reportService.CreateReport()
}

func (r *Runner) TagFile(filePath string) {
	//	TODO
	//	for block in file run GetBlameForFileLines and call tagger.CreateTagsForBlock

}
