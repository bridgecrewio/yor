package structure

import (
	"encoding/json"
	"fmt"
	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/minamijoyo/tfschema/tfschema"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type DockerfileParser struct {
	Labels                 map[string]string
	rootDir                string
	taggableResourcesCache map[string]bool
	providerToClientMap    sync.Map
	dockerFileExists       bool
	NewLabels              []TagItem
}

func (p *DockerfileParser) Name() string {
	return "Dockerfile"
}

func (p *DockerfileParser) ValidFile(_ string) bool {
	return true
}

func (p *DockerfileParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	var parsedBlocks []structure.IBlock
	var block structure.Block
	block.IsTaggable = false
	sArray, errReadFile := readFile(filePath)
	if errReadFile != nil {
		log.Fatal(errReadFile)
	}
	for _, line := range sArray {
		if strings.Contains(strings.ToLower(line), "label") {
			tag := generateTag(line)
			block.ExitingTags = append(block.ExitingTags, &tag)

		}
	}

	files, err := ioutil.ReadDir(p.rootDir)
	if err != nil {
		log.Fatal(err)
	}
	foundConfig := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == common.JSONFileType.Extension {
			foundConfig = true
			src, err := ioutil.ReadFile(file.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s because %s", filePath, err)
			}
			var input InputTags
			json.Unmarshal(src, &input.Input)
			for key, value := range input.Input {
				var tag tags.Tag
				tag.Key = key
				tag.SetValue(FormatString(value))
				block.NewTags = append(block.NewTags, &tag)
				block.IsTaggable = true

			}
		}
	}
	if !foundConfig {
		return nil, fmt.Errorf("Failed to find input JSON file for configureation")
	}

	parsedBlocks = append(parsedBlocks, &block)

	return parsedBlocks, nil

}

func FormatString(inputString string) string {
	chars := []rune(inputString)
	var newChars []rune
	for i := 0; i < len(chars); i++ {
		if string(chars[i]) != "\"" {
			newChars = append(newChars, chars[i])
		}
	}
	return string(newChars)

}

func (p *DockerfileParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	var tags []TagItem
	for _, item := range blocks {
		for _, oldTag := range item.GetExistingTags() {
			var tagItem TagItem
			tagItem.Key = oldTag.GetKey()
			tagItem.Value = oldTag.GetValue()
			tagItem.Exists = true
			tagItem.NeedsUpdating = false
			tagItem.Remove = existsInNewLabels(tagItem, item)
			tags = append(tags, tagItem)

		}
		for _, newTag := range item.GetNewTags() {
			var tagItem TagItem
			tagItem.Key = newTag.GetKey()
			tagItem.Value = newTag.GetValue()
			result := existsInOldLabels(tagItem.Key, tags)
			if result {
				tags = checkItemValue(tagItem, tags)
			} else {
				tagItem.Exists = false
				tagItem.NeedsUpdating = false
				tags = append(tags, tagItem)
			}

		}
	}

	sArray, errReadFile := readFile(readFilePath)
	if errReadFile != nil {
		log.Fatal(errReadFile)
	}
	fileChanged := false
	var newFile string
	for _, line := range sArray {
		if strings.Contains(strings.ToLower(line), "label") {
			tag := generateTag(line)
			for _, tagItem := range tags {
				if tagItem.Key == tag.Key {
					if tagItem.Value != tag.Value && tagItem.Remove == false {
						newLine := "LABEL " + tagItem.Key + "=" + "\"" + tagItem.Value + "\"\n"
						newFile = newFile + newLine
						fileChanged = true
					}
					if tagItem.Value == tag.Value && tagItem.Remove == false {
						newFile = newFile + line + "\n"
						fileChanged = true
					}
					if tagItem.Remove {

					}
				}
			}
		} else {
			newFile = newFile + line + "\n"
		}
	}
	for _, tagItem := range tags {
		if !tagItem.Exists {
			newLine := "LABEL " + tagItem.Key + "=" + "\"" + tagItem.Value + "\"\n"
			newFile = newFile + newLine
			fileChanged = true
		}
	}
	if fileChanged {

		var fileError = os.WriteFile(writeFilePath, []byte(newFile), 0644)
		if fileError != nil {
			var errorMessage = fmt.Errorf("failed to read file %s because %s", writeFilePath, fileError)
			log.Fatal(errorMessage)
		}
	}

	return nil
}

func existsInNewLabels(tag TagItem, block structure.IBlock) bool {
	newLabels := block.GetNewTags()
	for _, item := range newLabels {
		if tag.Key == item.GetKey() {
			return false
		}
	}
	return true
}

func checkItemValue(tag TagItem, list []TagItem) []TagItem {
	index := 0
	for _, item := range list {
		if item.Key == tag.Key {
			if item.Value != tag.Value {
				var updatedTag TagItem
				updatedTag.Key = tag.Key
				updatedTag.Value = tag.Value
				updatedTag.NeedsUpdating = true
				updatedTag.Exists = true
				list[index] = updatedTag
			}
		}
		index += 1
	}
	return list
}

func existsInOldLabels(key string, list []TagItem) bool {
	for _, item := range list {
		if item.Key == key {
			return true
		}
	}
	return false
}

func (p *DockerfileParser) GetSkippedDirs() []string {
	var ignoredDirs []string
	return ignoredDirs
}

func (p *DockerfileParser) Close() {
	logger.MuteOutputBlock(func() {
		p.providerToClientMap.Range(func(provider, iClient interface{}) bool {
			client := iClient.(tfschema.Client)
			client.Close()
			return true
		})
	})
}

func (p *DockerfileParser) GetSupportedFileExtensions() []string {
	return []string{common.DockerFileType.Extension}
}

func (p *DockerfileParser) Init(rootDir string, args map[string]string) {
	p.rootDir = rootDir
	//p.dockerFileExists = false
	//p.IngestFile()
	//p.ReadYamlConfiguration()
	//p.PatchDockerFile()

}

func generateTag(line string) tags.Tag {
	reg := regexp.MustCompile(`[0-9a-zA-Z\"]+=[0-9a-zA-Z@.\"\s]+`)
	result := reg.FindString(line)
	keyValue := strings.Split(result, "=")
	if len(keyValue) < 2 {
		log.Fatal("Line in Dockerfile could not be parsed correctly.")
	}
	tag := tags.Tag{
		Key:   FormatString(keyValue[0]),
		Value: FormatString(keyValue[1]),
	}
	return tag
}

func readFile(filePath string) ([]string, error) {
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s because %s", filePath, err)
	}
	sArray := strings.Split(string(src), "\n")
	return sArray, nil
}

type TagItem struct {
	Key           string
	Value         string
	NeedsUpdating bool
	Exists        bool
	Remove        bool
}

type InputTags struct {
	Input map[string]string
}
