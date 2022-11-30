package code2cloud

import (
	"fmt"

	"github.com/bridgecrewio/yor/src/common/tagging/tags"

	"github.com/google/uuid"
)

type YorTraceTag struct {
	tags.Tag
}

func (t *YorTraceTag) Init() {
	t.Key = tags.YorTraceTagKey
}

func (t *YorTraceTag) CalculateValue(_ interface{}) (tags.ITag, error) {
	uuidv4, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create a new uuidv4")
	}
	return &tags.Tag{Key: t.Key, Value: uuidv4.String()}, nil
}

func (t *YorTraceTag) GetDescription() string {
	return "A UUID tag that allows easily finding the root IaC config of the resource"
}
