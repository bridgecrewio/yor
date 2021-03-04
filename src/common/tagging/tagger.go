package tagging

import (
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"regexp"
	"strings"
)

type Tagger struct {
	tags        []tags.ITag
	SkippedTags []string
}

var IgnoredDirs = []string{".git", ".DS_Store", ".idea"}

type ITagger interface {
	InitTagger(path string, skippedTags []string)
	CreateTagsForBlock(block structure.IBlock)
	GetTags() []tags.ITag
}

func (t *Tagger) GetSkippedDirs() []string {
	return IgnoredDirs
}

func (t *Tagger) SetTags(tags []tags.ITag) {
	for _, tag := range tags {
		tag.Init()
		if !t.IsTagSkipped(tag) {
			t.tags = append(t.tags, tag)
		}
	}
}

func (t *Tagger) GetTags() []tags.ITag {
	return t.tags
}

func (t *Tagger) IsTagSkipped(tag tags.ITag) bool {
	for _, st := range t.SkippedTags {
		stRegex := strings.ReplaceAll(st, "*", ".*")
		if match, err := regexp.Match(stRegex, []byte(tag.GetKey())); match || err != nil {
			logger.Info(fmt.Sprintf("Skipping %v due to skip-tag constraint %v", tag.GetKey(), st))
			return true
		}
	}
	return false
}
