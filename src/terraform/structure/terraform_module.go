package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/moduledeps"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/mitchellh/cli"
)

const PluginsOutputDir = ".plugins"

var SkippedProviders = []string{"null", "random", "tls"}

type TerraformModule struct {
	tfModule            *tfconfig.Module
	rootDir             string
	ProvidersInstallDir string
}

func NewTerraformModule(rootDir string) *TerraformModule {
	tfModule, diagnostics := tfconfig.LoadModule(rootDir)
	if diagnostics != nil && diagnostics.HasErrors() {
		logger.Error(diagnostics.Error())
		return nil
	}
	terraformModule := &TerraformModule{tfModule: tfModule, rootDir: rootDir}
	// download terraform plugin into local folder if it doesn't exist
	pwd, _ := os.Getwd()
	terraformModule.ProvidersInstallDir = path.Join(pwd, PluginsOutputDir)
	terraformModule.InitProvider()

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
		if common.InSlice(SkippedProviders, provider) {
			continue
		}
		if providerExists(t.ProvidersInstallDir, provider) {
			return
		}
		pty := addrs.NewLegacyProvider(provider)
		logger.MuteLogging()
		//if provider == "google" {
		//	setGoogleVersion, err := discovery.ConstraintStr("3.65.0").Parse()
		//	if err != nil {
		//		logger.Error(fmt.Sprintf("failed to parse google version str because of errors %s", err))
		//	}
		//	constraints.Versions.Append(setGoogleVersion)
		//}
		_, diagnostics, err := providerInstaller.Get(pty, constraints.Versions)
		logger.UnmuteLogging()
		if diagnostics != nil && diagnostics.HasErrors() {
			logger.Error(fmt.Sprintf("failed to install provider for directory %s because of errors %s", t.rootDir, diagnostics.Err()))
		}
		if err != nil {
			logger.Error(fmt.Sprintf("failed to install provider for directory %s because of errors %s", t.rootDir, err))
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
		childModuleDir := path.Join(t.rootDir, moduleCall.Source)
		childModule := NewTerraformModule(childModuleDir)
		childModulesDirectories := childModule.GetModulesDirectories()
		for _, childDirPath := range childModulesDirectories {
			if _, err := os.Stat(childDirPath); !os.IsNotExist(err) && !common.InSlice(modulesDirectories, childDirPath) {
				// if directory exists (local module) and modulesDirectories doesn't contain it yet, add it
				modulesDirectories = append(modulesDirectories, childDirPath)
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
		if strings.HasPrefix(moduleCall.Source, "git::") {
			logger.Info("Skipping remote module", moduleCall.Source)
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
