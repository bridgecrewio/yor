package tags

import (
	"fmt"

	"github.com/google/uuid"
)

type YorTraceTag struct {
	Tag
}

func (t *YorTraceTag) Init() {
	t.Key = "yor_trace"
}

func (t *YorTraceTag) CalculateValue(_ interface{}) error {
	uuidv4, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to create a new uuidv4")
	}
	t.Value = uuidv4.String()
	return nil
}
