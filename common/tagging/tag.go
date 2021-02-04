package tagging

type Tag struct {
	Key   string
	Value interface{}
}

type ITag interface {
	Init() ITag
	CalculateValue(data interface{}) error
	GetTag() map[string]interface{}
}

func Init(key string, value interface{}) ITag {
	return &Tag{
		Key:   key,
		Value: value,
	}
}

func (t *Tag) Init() ITag {
	return t
}

func (t *Tag) CalculateValue(data interface{}) error {
	return nil
}

func (t *Tag) GetTag() map[string]interface{} {
	return map[string]interface{}{t.Key: t.Value}
}
