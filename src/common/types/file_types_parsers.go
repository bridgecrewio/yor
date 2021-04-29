package types

import (
	"bridgecrewio/yor/src/common/structure"
)


type YamlParser struct {
	RootDir              string
	FileToResourcesLines map[string]structure.Lines
}

