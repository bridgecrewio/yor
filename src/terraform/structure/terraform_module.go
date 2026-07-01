package structure

import (
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/utils"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

const PluginsOutputDir = ".yor_plugins"

var SkippedProviders = []string{"null", "random", "tls", "local", "sops"}
var RegistryModuleRegex = regexp.MustCompile("^((?P<MODULE_HOSTNAME>[^/]+)/)?(?P<MODULE_NAMESPACE>[^/]+)/(?P<MODULE_NAME>[^/]+)/(?P<PROVIDER>[a-z]+)")

type TerraformModule struct {
	tfModule *tfconfig.Module
	rootDir  string
}

func NewTerraformModule(rootDir string) *TerraformModule {
	tfModule, diagnostics := tfconfig.LoadModule(rootDir)
	if diagnostics != nil && diagnostics.HasErrors() {
		logger.Warning(diagnostics.Error())
		return nil
	}
	// Provider plugins are no longer installed/launched: taggability is resolved from the
	// static TfTaggableResourceTypes list, so yor does not need the legacy terraform
	// provider-install machinery (which pulled the vulnerable terraform v0.12.31 chain).
	return &TerraformModule{tfModule: tfModule, rootDir: rootDir}
}

func (t *TerraformModule) GetModulesDirectories() []string {
	modulesDirectories := []string{t.rootDir}

	for _, moduleCall := range t.tfModule.ModuleCalls {
		if !isRemoteModule(moduleCall.Source) && !isTerraformRegistryModule(moduleCall.Source) {
			childModuleDir := path.Join(t.rootDir, moduleCall.Source)
			childModule := NewTerraformModule(childModuleDir)
			childModulesDirectories := childModule.GetModulesDirectories()
			for _, childDirPath := range childModulesDirectories {
				if _, err := os.Stat(childDirPath); !os.IsNotExist(err) && !utils.InSlice(modulesDirectories, childDirPath) {
					// if directory exists (local module) and modulesDirectories doesn't contain it yet, add it
					modulesDirectories = append(modulesDirectories, childDirPath)
				}
			}
		}
	}

	return modulesDirectories
}

func isRemoteModule(s string) bool {
	// Taken from https://www.terraform.io/docs/language/modules/sources.html
	return strings.HasPrefix(s, "git::") || strings.HasPrefix(s, "hg::") || strings.HasPrefix(s, "s3::") || strings.HasPrefix(s, "gcs::") ||
		strings.HasPrefix(s, "github.com/") || strings.HasPrefix(s, "bitbucket.org/") || strings.HasPrefix(s, "app.terraform.io/") ||
		strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "git@")
}

func isTerraformRegistryModule(source string) bool {
	matches := utils.FindSubMatchByGroup(RegistryModuleRegex, source)
	if matches == nil {
		return false
	}
	if provider, ok := matches["PROVIDER"]; ok {
		if _, okTag := ProviderToTagAttribute[provider]; okTag {
			return true
		}
	}
	return false
}
