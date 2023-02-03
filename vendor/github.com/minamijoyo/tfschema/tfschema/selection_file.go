package tfschema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// selectionFile represents a lock file in Terraform v0.13.
// An example is as follows:
// ```.terraform/plugins/selections.json
// {
//   "registry.terraform.io/hashicorp/aws": {
//     "hash": "h1:bVpVG796X94WeMxRcRNq+YmHVQkbaWYCsR906VwgJxE=",
//     "version": "2.67.0"
//   }
// }
// ```
//
// It includes selected the fully qualified provider name and version.
// It is implemented in github.com/hashicorp/terraform/internal/getproviders.lockFile,
// but we cannot import the internal package.
// So we peek the file and decode information only we need here.
type selectionFile struct {
	path string
}

// selectionFileEntry represents an entry in SelectionFile.
type selectionFileEntry struct {
	// Version is a selected version.
	Version string `json:"version"`
}

func newSelectionFile(path string) *selectionFile {
	return &selectionFile{
		path: path,
	}
}

// pluginDirs peeks a selection file and returns a slice of plugin directories.
func (f *selectionFile) pluginDirs() ([]string, error) {
	buf, err := ioutil.ReadFile(f.path)

	if err != nil {
		if os.IsNotExist(err) {
			// not found, just ignore it
			log.Printf("[DEBUG] selection file not found: %s", f.path)
			return []string{}, nil
		}
		return nil, err
	}

	var entries map[string]selectionFileEntry
	err = json.Unmarshal(buf, &entries)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse selection file %s: %s", f.path, err)
	}

	dirs := []string{}
	base := filepath.Dir(f.path)
	arch := runtime.GOOS + "_" + runtime.GOARCH
	for k, v := range entries {
		// build a path for plugin dir such as the following:
		// .terraform/plugins/registry.terraform.io/hashicorp/aws/2.67.0/darwin_amd64
		dir := filepath.Join(base, k, v.Version, arch)
		dirs = append(dirs, dir)
	}

	log.Printf("[DEBUG] selection file found: %s: %v", f.path, dirs)
	return dirs, nil
}
