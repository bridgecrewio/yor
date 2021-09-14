package types

import (
	"sync"
)

type YamlParser struct {
	RootDir              string
	FileToResourcesLines sync.Map
}

type JSONParser struct {
	RootDir              string
	FileToBracketMapping sync.Map
}
