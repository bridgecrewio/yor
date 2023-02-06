package tags

import (
	"fmt"
	"regexp"
)

type Tag struct {
	Key   string
	Value string
}

const YorTraceTagKey = "yor_trace"
const GitFileTagKey = "git_file"
const GitModifiersTagKey = "git_modifiers"
const GitLastModifiedAtTagKey = "git_last_modified_at"
const GitLastModifiedByTagKey = "git_last_modified_by"
const GitRepoTagKey = "git_repo"

type ITag interface {
	Init()
	CalculateValue(data interface{}) (ITag, error)
	GetKey() string
	SetValue(val string)
	GetValue() string
	GetPriority() int
	GetDescription() string
	SetTagPrefix(tagPrefix string)
}

type TagDiff struct {
	Key       string
	PrevValue string
	NewValue  string
}

func Init(key string, value string) ITag {
	return &Tag{
		Key:   key,
		Value: value,
	}
}

func (t *Tag) Init() {}

func (t *Tag) SetTagPrefix(tagPrefix string) {
	t.Key = fmt.Sprintf("%s%s", tagPrefix, t.Key)
}

func (t *Tag) GetPriority() int {
	return 0
}

func (t *Tag) CalculateValue(_ interface{}) (ITag, error) {
	return &Tag{
		Key:   t.Key,
		Value: t.Value,
	}, nil
}

func (t *Tag) GetDescription() string {
	return "Abstract tag class"
}

func (t *Tag) GetKey() string {
	return t.Key
}

func (t *Tag) SetValue(val string) {
	t.Value = val
}

func (t *Tag) GetValue() string {
	return t.Value
}

// IsTagKeyMatch Try to match the tag's key name with a potentially quoted string
func IsTagKeyMatch(tag ITag, keyName string) bool {
	match, _ := regexp.Match(fmt.Sprintf(`\b"?%s"?\b`, regexp.QuoteMeta(keyName)), []byte(tag.GetKey()))
	return match
}
