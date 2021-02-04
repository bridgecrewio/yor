package runner

import (
	"bridgecrewio/yor/reports"
	"bridgecrewio/yor/tagging"
	"bridgecrewio/yor/tagging/taggers"
)

type Runner struct {
	taggers []taggers.ITagger
}

func NewRunner(taggerTypes []tagging.TaggerType, extraTags []tagging.ITag) *Runner {
	// TODO
	return nil
}

func initTaggers(taggerTypes []tagging.TaggerType, extraTags []tagging.ITag) {
	// TODO: create a new Tagger instance and send the extraTags as param
}

func (r *Runner) TagDirectory(dir string) interface{} {
	//	TODO - return Report's result from this method
	reportService := reports.ReportService{}

	return reportService.CreateReport()
}

func (r *Runner) TagFile(filePath string) {
	//	TODO
	//	for block in file run GetBlameForFileLines and call tagger.CreateTagsForBlock

}
