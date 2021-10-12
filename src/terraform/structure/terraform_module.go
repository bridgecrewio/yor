package structure

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/utils"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/moduledeps"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/mitchellh/cli"
)

const PluginsOutputDir = ".yor_plugins"

var SkippedProviders = []string{"null", "random", "tls", "local"}
var RegistryModuleRegex = regexp.MustCompile("^(?P<MODULE_WRITER>[^/]+)/(?P<MODULE_NAME>[^/]+)/(?P<PROVIDER>[a-z]+)")

type TerraformModule struct {
	tfModule            *tfconfig.Module
	rootDir             string
	ProvidersInstallDir string
}

func NewTerraformModule(rootDir string) *TerraformModule {
	tfModule, diagnostics := tfconfig.LoadModule(rootDir)
	if diagnostics != nil && diagnostics.HasErrors() {
		logger.Warning(diagnostics.Error())
		return nil
	}
	terraformModule := &TerraformModule{tfModule: tfModule, rootDir: rootDir}
	if strings.ToUpper(os.Getenv("YOR_SKIP_PROVIDER_DOWNLOAD")) != "TRUE" {
		// download terraform plugin into local folder if it doesn't exist
		homeDir, _ := os.UserHomeDir()
		terraformModule.ProvidersInstallDir = path.Join(homeDir, PluginsOutputDir)
		terraformModule.InitProvider()
	}

	return terraformModule
}

func (t *TerraformModule) InitProvider() {
	moduleDependencies := getProviderDependencies(t.tfModule)
	providers := moduleDependencies.AllPluginRequirements()
	providerInstaller := &discovery.ProviderInstaller{
		Dir:                   t.ProvidersInstallDir,
		PluginProtocolVersion: discovery.PluginInstallProtocolVersion,
		SkipVerify:            false,
		Ui:                    &cli.MockUi{},
	}
	for provider, constraints := range providers {
		if utils.InSlice(SkippedProviders, provider) {
			continue
		}
		if providerExists(t.ProvidersInstallDir, provider) {
			return
		}
		pty := addrs.NewLegacyProvider(provider)
		logger.MuteLogging()
		_, diagnostics, err := providerInstaller.Get(pty, constraints.Versions)
		logger.UnmuteLogging()
		if (diagnostics != nil && diagnostics.HasErrors()) || err != nil {
			errMsg := diagnostics.Err()
			if errMsg == nil {
				errMsg = err
			}
			logger.Warning(fmt.Sprintf("failed to install provider \"%v\" for directory %s because of errors %s", provider, t.rootDir, errMsg))
		}
	}
}

func providerExists(providersInstallDir string, provider string) bool {
	fileInfo, err := ioutil.ReadDir(providersInstallDir)
	if err != nil {
		return false
	}
	for _, file := range fileInfo {
		if strings.Contains(file.Name(), provider) && strings.Contains(file.Name(), "provider") {
			return true
		}
	}

	return false
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

func getProviderDependencies(tfModule *tfconfig.Module) *moduledeps.Module {
	moduleDependencies := &moduledeps.Module{}
	providers := make(moduledeps.Providers)

	for name, requirement := range tfModule.RequiredProviders {
		var constraints version.Constraints
		for _, reqStr := range requirement.VersionConstraints {
			if reqStr != "" {
				constraint, err := version.NewConstraint(reqStr)
				if err != nil {
					logger.Warning(fmt.Sprintf("Invalid version constraint %q for provider %s.", reqStr, name))
					continue
				}
				constraints = append(constraints, constraint...)
			}
		}

		inst := moduledeps.ProviderInstance(name)
		providers[inst] = moduledeps.ProviderDependency{
			Constraints: discovery.NewConstraints(constraints),
			Reason:      moduledeps.ProviderDependencyExplicit,
		}
	}

	for name := range ProviderToTagAttribute {
		inst := moduledeps.ProviderInstance(name)
		if _, ok := providers[inst]; !ok {
			providers[inst] = moduledeps.ProviderDependency{
				Constraints: discovery.Constraints{},
				Reason:      moduledeps.ProviderDependencyImplicit,
			}
		}
	}
	moduleDependencies.Providers = providers

	for _, moduleCall := range tfModule.ModuleCalls {
		if isRemoteModule(moduleCall.Source) || isTerraformRegistryModule(moduleCall.Source) {
			logger.Info("Skipping remote git module", moduleCall.Source)
			continue
		}
		childModulePath := path.Join(tfModule.Path, moduleCall.Source)
		tfChildModule, diagnostics := tfconfig.LoadModule(childModulePath)
		if diagnostics != nil && diagnostics.HasErrors() {
			hclErrors := diagnostics.Error()
			logger.Warning(fmt.Sprintf("failed to parse hcl module in directory %s because of errors %s", path.Join(childModulePath, moduleCall.Source), hclErrors))
		} else {
			child := getProviderDependencies(tfChildModule)
			moduleDependencies.Children = append(moduleDependencies.Children, child)
		}
	}

	return moduleDependencies
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
