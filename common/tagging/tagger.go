package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/tagging/tags"
	"log"
)

type Tagger struct {
	tags             []tags.ITag
	TagAttributeName string
}

type ITagger interface {
	IsBlockTaggable(block interface{}) bool
	CreateTagsForBlock(block interface{}, gitBlame *git_service.GitBlame) // this method should call IBlock.AddNewTags
}

func (t *Tagger) InitTags(extraTags []tags.ITag) {
	// TODO: initialize Tagger.tags with all the tags in `TagTypes` (use their Init method) and add the extra tags
	for _, tagType := range tags.TagTypes {
		log.Print(tagType)
	}

	for _, extraTagType := range extraTags {
		log.Print(extraTagType)
	}
}
