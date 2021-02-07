package tags

type Tag struct {
	Key   string
	Value string
}

var TagTypes = []ITag{
	&GitOrgTag{},
	&GitRepoTag{},
	&GitFileTag{},
	&GitCommitTag{},
	&GitModifiersTag{},
	&GitLastModifiedAtTag{},
	&GitLastModifiedByTag{},
}

type ITag interface {
	Init() ITag
	CalculateValue(data interface{}) error
	GetTag() map[string]string
}

func Init(key string, value string) ITag {
	return &Tag{
		Key:   key,
		Value: value,
	}
}

func (t *Tag) Init() ITag {
	return t
}

func (t *Tag) CalculateValue(_ interface{}) error {
	return nil
}

func (t *Tag) GetTag() map[string]string {
	return map[string]string{t.Key: t.Value}
}
