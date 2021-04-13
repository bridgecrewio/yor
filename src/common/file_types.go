package common

type FileType struct {
	Extension  string
	FileFormat string
}

var YamlFileType = FileType{Extension: ".yaml", FileFormat: "yaml"}
var YmlFileType = FileType{Extension: ".yml", FileFormat: "yml"}
var JSONFileType = FileType{Extension: ".json", FileFormat: "json"}
var CFTFileType = FileType{Extension: ".template", FileFormat: "template"}
var TfFileType = FileType{Extension: ".tf", FileFormat: "tf"}
