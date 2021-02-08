package structure

type TerrraformParser struct {
}

func (p *TerrraformParser) ParseFile(filePath string) ([]*TerraformBlock, error) {
	// TODO
	return nil, nil
}

func (p *TerrraformParser) WriteFile(filePath string, blocks []*TerraformBlock) error {
	// TODO
	return nil
}
