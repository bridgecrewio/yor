package types

import (
	"github.com/bridgecrewio/yor/src/common/json"
	"github.com/bridgecrewio/yor/src/common/structure"
)

type YamlParser struct {
	RootDir              string
	FileToResourcesLines map[string]structure.Lines
}

type JSONParser struct {
	RootDir              string
	FileToBracketMapping map[string]map[int]json.BracketPair
}
