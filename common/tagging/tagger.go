package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
)

var TagTypes = []ITag{&tags.GitCommitTag{}}

type Tagger struct {
	tags             []ITag
	TagAttributeName string
}

type ITagger interface {
	IsBlockTaggable(block structure.IBlock) bool
	CreateTagsForBlock(block structure.IBlock, gitBlame *git_service.GitBlame) // this method should call Block.AddNewTags
}

func (t *Tagger) InitTags(extraTags []ITag) {
	// TODO: initialize Tagger.tags with all the tags in `TagTypes` (use their Init method) and add the extra tags
}
