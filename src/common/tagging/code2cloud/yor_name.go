package code2cloud

import (
	"fmt"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"reflect"
)

type YorNameTag struct {
	tags.Tag
}

func (t *YorNameTag) Init() {
	t.Key = tags.YorNameTagKey
}

func (t *YorNameTag) CalculateValue(block interface{}) (tags.ITag, error) {
	blockVal, ok := block.(structure.IBlock)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to IBlock, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(block))
	}

	return &tags.Tag{Key: t.Key, Value: blockVal.GetResourceName()}, nil
}

func (t *YorNameTag) GetDescription() string {
	return "A tag that states the resource name in the IaC config file"
}
