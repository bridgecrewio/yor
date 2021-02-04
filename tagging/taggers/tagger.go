package taggers

import (
	"bridgecrewio/yor/git_service"
	"bridgecrewio/yor/structure"
	"bridgecrewio/yor/tagging"
	"bridgecrewio/yor/tagging/tags"
)

var TagTypes = []tagging.ITag{&tags.GitCommitTag{}}

type Tagger struct {
	tags             []*tagging.Tag
	TagAttributeName string
}

type ITagger interface {
	IsBlockTaggable(block *structure.Block) bool
	CreateTagsForBlock(block *structure.Block, gitBlame *git_service.GitBlame) // this method should call Block.AddNewTags
}

func (t *Tagger) InitTags(extraTags []tagging.ITag) {
	// TODO: initialize Tagger.tags with all the tags in `TagTypes` (use their Init method) and add the extra tags
}
