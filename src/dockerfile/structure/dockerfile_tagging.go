package structure

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (p *DockerfileParser) IngestFile() {
	//  Read in the file with the file name of either "Dockerfile" or "dockerfile"
	//  All other files will not work until we discover a better way of verifying a dockerfile

	filePtr := p.GetDockerfile(p.rootDir)
	if filePtr == nil {
		log.Fatal("Error finding dockerfile.")
	} else {
		p.getDockerLabel(filePtr)
	}
}

func (p *DockerfileParser) PatchDockerFile() {
	filePtr := p.GetDockerfile(p.rootDir)
	file, fileResults := p.normalizeDockerLabel(filePtr)
	if fileResults {
		d1 := []byte(file)
		err := os.WriteFile(p.rootDir+"/Dockerfile", d1, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (p *DockerfileParser) GetDockerfile(directory string) *string {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == "Dockerfile" || file.Name() == "dockerfile" {
			p.dockerFileExists = true
			fileData, fileError := os.ReadFile(file.Name())
			if fileError != nil {
				log.Fatal(fileData)
			}
			filePointer := string(fileData)
			return &filePointer
		}
	}
	return nil
}

func (p *DockerfileParser) normalizeDockerLabel(filePtr *string) (string, bool) {
	dockerFile := *filePtr
	sArray := strings.Split(dockerFile, "\n")
	fileChanged := false
	var file string
	for _, line := range sArray {
		if strings.Contains(strings.ToLower(line), "label") {
			reg := regexp.MustCompile(`\w+=[a-zA-Z@.\"\s]+`)
			result := reg.FindString(line)
			keyValue := strings.Split(result, "=")
			tag := p.getTag(keyValue[0])
			if tag != nil && tag.NeedsUpdating {
				line = "LABEL " + tag.Key + "=\"" + tag.Value + "\"\n"
				file = file + line
				fileChanged = true
			}
		}
		if !strings.Contains(strings.ToLower(line), "label") {
			file = file + line + "\n"
		}
	}
	for _, item := range p.NewLabels {
		if !item.Exists {
			line := "LABEL " + item.Key + "=\"" + item.Value + "\"\n"
			file = file + line
			fileChanged = true
		}
	}
	return file, fileChanged
}

func (p *DockerfileParser) getDockerLabel(filePtr *string) {
	dockerFile := *filePtr
	sArray := strings.Split(dockerFile, "\n")
	m := make(map[string]string)
	for _, line := range sArray {
		if strings.Contains(strings.ToLower(line), "label") {
			reg := regexp.MustCompile(`\w+=[a-zA-Z@.\"\s]+`)
			result := reg.FindString(line)
			keyValue := strings.Split(result, "=")

			m[keyValue[0]] = keyValue[1]
		}
	}
	p.Labels = m
}

func (p *DockerfileParser) ReadYamlConfiguration() {
	files, err := ioutil.ReadDir(p.rootDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		m := make(map[string]string)
		extension := filepath.Ext(file.Name())
		if extension == ".yml" {
			fileData, fileError := os.ReadFile(file.Name())
			if fileError != nil {
				log.Fatal(fileData)
			}
			yaml.Unmarshal(fileData, &m)
			p.checkTags(m)
		}
	}
}

func (p *DockerfileParser) checkTags(configFile map[string]string) {
	var tags []TagItem
	for key, value := range configFile {
		tagItem := p.existsInMap(key, value)
		if tagItem != nil {
			tags = append(tags, *tagItem)
		}
	}
	p.NewLabels = tags
}

func (p *DockerfileParser) getTag(key string) *TagItem {
	for _, item := range p.NewLabels {
		if item.Key == key {
			return &item
		}
	}
	return nil
}

func (p *DockerfileParser) existsInMap(key string, value string) *TagItem {
	var tag TagItem
	tag.Key = key
	tag.Value = value
	tag.Exists = false
	tag.NeedsUpdating = false
	for sourceKey, sourceValue := range p.Labels {
		if sourceKey == key {
			tag.Exists = true
			if sourceValue != value {
				tag.NeedsUpdating = true
				return &tag
			} else {
				tag.NeedsUpdating = false
				return &tag
			}
		}

	}
	return &tag
}
