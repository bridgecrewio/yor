package yaml

import "github.com/bridgecrewio/yor/src/common/structure"

type IYamlBlock interface {
	structure.IBlock
	UpdateTags()
}
