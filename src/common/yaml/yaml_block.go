package yaml

import "bridgecrewio/yor/src/common/structure"

type IYamlBlock interface {
	structure.IBlock
	UpdateTags()
}
