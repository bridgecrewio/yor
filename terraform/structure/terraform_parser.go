package structure

import "bridgecrewio/yor/common/structure"

type TerrraformParser struct {
}

func (p *TerrraformParser) ParseFile(filePath string) ([]*structure.Block, error) {
	// TODO
	return nil, nil
}

func (p *TerrraformParser) WriteFile(filePath string, blocks []*structure.Block) error {
	// TODO
	return nil
}
