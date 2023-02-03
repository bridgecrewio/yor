package tfschema

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

// lockFile represents a lock file in Terraform v0.14+.
// An example is as follows:
// ```.terraform.lock.hcl
// # This file is maintained automatically by "terraform init".
// # Manual edits may be lost in future updates.
//
// provider "registry.terraform.io/hashicorp/aws" {
//   version     = "3.17.0"
//   constraints = "3.17.0"
//   hashes = [
//     "h1:BLK4zgpn2O4ojSZhCtAXzqsvm7BymqSxrtcXLp/J/yA=",
//     "zh:047e22b2e02d57fb1d945d52c4bd062f50a657865b6a8d21f96ba55ef8a474e5",
//     "zh:30135eb8eed7f4c135c504889cf66a2a586671783122a2953856b7040074cba9",
//     "zh:557a57801e5c004178caf8c9d94282b17bfab0c6f8fe68e414a4f650b441fc3c",
//     "zh:6583fbc0159a7835235d7e7b7ce08b5cc51514dc5bc8276c49b3e14a79a7ddea",
//     "zh:6dc01ed7ea928caf1fd5ace9a25e219b57642925d8784046f994e241adad45d2",
//     "zh:9c8a1942c21dcb78992a39d2ca54a08d839cc936a68341fb5443983ebd17cad0",
//     "zh:c71ca4ccb2c228bfb1aa5e127a54872c739fe54c477b4c545df2e74e77ef3bf4",
//     "zh:c94ed23ed6770721262d7bc342467c02524bb0380ece0122f224e30a6b0d23de",
//     "zh:e3b71cb43d0602ee1fcf520a7bd7a2a9232b751138e6ef5cf6837caa26c63173",
//     "zh:f5143912cd4588bc3a19a98a7fd7979202b9c55130c77b2697cf64bfea76c61c",
//   ]
// }
// ```
//
// It includes selected the fully qualified provider name and version.
// It is implemented in github.com/hashicorp/terraform/internal/internal/depsfile.ProviderLock,
// but we cannot import the internal package.
// So we peek the file and decode information only we need here.
type lockFile struct {
	path string
}

// newLockFile returns a new lockFile instance.
func newLockFile(path string) *lockFile {
	return &lockFile{
		path: path,
	}
}

// pluginDirs peeks a lock file and returns a slice of plugin directories.
func (f *lockFile) pluginDirs() ([]string, error) {
	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		// not found, just ignore it
		log.Printf("[DEBUG] lock file not found: %s", f.path)
		return []string{}, nil
	}

	lock, err := loadLockFile(f.path)
	if err != nil {
		return nil, err
	}

	dirs := []string{}
	base := filepath.Dir(f.path)
	arch := runtime.GOOS + "_" + runtime.GOARCH
	for _, p := range lock.Providers {
		// build a path for plugin dir such as the following:
		// .terraform/providers/registry.terraform.io/hashicorp/aws/3.17.0/darwin_amd64
		dir := filepath.Join(base, ".terraform", "providers", p.Address, p.Version, arch)
		dirs = append(dirs, dir)
	}

	log.Printf("[DEBUG] lock file found: %s: %v", f.path, dirs)
	return dirs, nil
}

// Lock represents a lock file written in HCL.
type Lock struct {
	// A list of providers.
	Providers []Provider `hcl:"provider,block"`
	// The rest of body we don't need.
	Remain hcl.Body `hcl:",remain"`
}

// Provider represents a provider block in HCL.
type Provider struct {
	// A fully qualified provider name. (e.g. registry.terraform.io/hashicorp/aws)
	Address string `hcl:"address,label"`
	// A selected version.
	Version string `hcl:"version"`
	// The rest of body we don't need.
	Remain hcl.Body `hcl:",remain"`
}

// loadLockFile loads and parses a given lock file.
func loadLockFile(path string) (*Lock, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parseLockFile(path, source)
}

// parseLockFile parses a given lock file.
func parseLockFile(path string, source []byte) (*Lock, error) {
	var lock Lock
	err := hclsimple.Decode(path, source, nil, &lock)
	if err != nil {
		return nil, fmt.Errorf("failed to decode a lock file: %s, err: %s", path, err)
	}

	return &lock, nil
}
