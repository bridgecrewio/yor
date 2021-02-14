package main

type CheckovTag struct {
	Key   string
	Value string
}

func (t *CheckovTag) Init() {
	t.Key = "yor_checkov"
}

func (t *CheckovTag) CalculateValue(data interface{}) error {
	t.Value = "checkov"
	return nil
}

func (t *CheckovTag) GetKey() string {
	return t.Key
}

func (t *CheckovTag) GetValue() string {
	return t.Value
}

func (t *CheckovTag) GetPriority() int {
	return 1
}
