package tags

import (
	"github.com/google/uuid"
)

type BcTraceTag struct {
	Tag
}

func (t *BcTraceTag) Init() ITag {
	t.Key = "BC_TRACE"
	return t
}

func (t *BcTraceTag) CalculateValue(_ interface{}) error {
	t.Value = uuid.NewString()
	return nil
}
