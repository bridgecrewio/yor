package main

type TerragoatTag struct {
	Key   string
	Value string
}

func (t *TerragoatTag) Init() {
	t.Key = "yor_terragoat"
}

func (t *TerragoatTag) CalculateValue(data interface{}) error {
	t.Value = "terragoat"
	return nil
}

func (t *TerragoatTag) GetKey() string {
	return t.Key
}

func (t *TerragoatTag) GetValue() string {
	return t.Value
}

func (t *TerragoatTag) GetPriority() int {
	return -1
}
