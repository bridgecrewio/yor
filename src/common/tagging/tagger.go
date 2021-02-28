package tagging

import (
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
)

type Tagger struct {
	Tags []tags.ITag
}

var IgnoredDirs = []string{".git", ".DS_Store", ".idea"}

type ITagger interface {
	InitTagger(path string)
	CreateTagsForBlock(block structure.IBlock)
	GetDescription() string
}

func (t *Tagger) GetSkippedDirs() []string {
	return IgnoredDirs
}
