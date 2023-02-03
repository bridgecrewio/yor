package tfschema

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/mitchellh/go-homedir"
)

// Client represents a set of methods required to get a type definition of
// schema from Terraform providers.
// Terraform v0.12+ has a different provider interface from v0.11.
// This is a compatibility layer for Terraform v0.11/v0.12+.
type Client interface {
	// GetProviderSchema returns a type definiton of provider schema.
	GetProviderSchema() (*Block, error)

	// GetResourceTypeSchema returns a type definiton of resource type.
	GetResourceTypeSchema(resourceType string) (*Block, error)

	// GetDataSourceSchema returns a type definiton of data source.
	GetDataSourceSchema(dataSource string) (*Block, error)

	// ResourceTypes returns a list of resource types.
	ResourceTypes() ([]string, error)

	// DataSources returns a list of data sources.
	DataSources() ([]string, error)

	// Close closes a connection and kills a process of the plugin.
	Close()
}

// Option is an options struct for extra options for NewClient
type Option struct {
	RootDir string
	Logger  hclog.Logger
}

// NewClient creates a new Client instance.
func NewClient(providerName string, options Option) (Client, error) {
	// First, try to connect by GRPC protocol (version 5)
	log.Println("[DEBUG] try to connect by GRPC protocol (version 5)")
	client, err := NewGRPCClient(providerName, options)
	if err == nil {
		return client, nil
	}

	// If failed, try to connect by NetRPC protocol (version 4)
	// plugin.ClientConfig.AllowedProtocols has a protocol negotiation feature,
	// but it doesn't seems to work with old providers.
	// We guess it is for Terraform v0.11 to connect to the latest provider.
	// So we implement our own fallback logic here.
	log.Println("[DEBUG] try to connect by NetRPC protocol (version 4)")
	client, err = NewNetRPCClient(providerName, options)
	if err != nil {
		return nil, fmt.Errorf("Failed to NewClient: %s", err)
	}

	return client, nil
}

// findPlugin finds a plugin with the name specified in the arguments.
func findPlugin(pluginType string, pluginName string, rootDir string) (*discovery.PluginMeta, error) {
	dirs, err := pluginDirs(rootDir)
	if err != nil {
		return nil, err
	}

	pluginMetaSet := discovery.FindPlugins(pluginType, dirs).WithName(pluginName)

	// if pluginMetaSet doesn't have any pluginMeta, pluginMetaSet.Newest() will call panic.
	// so check it here.
	if pluginMetaSet.Count() > 0 {
		ret := pluginMetaSet.Newest()
		return &ret, nil
	}

	return nil, fmt.Errorf("Failed to find plugin: %s. Plugin binary was not found in any of the following directories: [%s]", pluginName, strings.Join(dirs, ", "))
}

// pluginDirs returns a list of directories to find plugin.
// It finds plugins installed by terraform init.
// Note that it doesn't have exactly the same behavior of Terraform
// because of some reasons:
// - Support multiple Terraform versions
// - Can't import internal packages of Terraform and it's too complicated to support
// - For debug
// For more details, read inline comments.
func pluginDirs(rootDir string) ([]string, error) {
	dirs := []string{}

	// current directory
	dirs = append(dirs, rootDir)

	// same directory as this executable (not terraform)
	exePath, err := os.Executable()
	if err != nil {
		return []string{}, fmt.Errorf("Failed to get executable path: %s", err)
	}
	dirs = append(dirs, filepath.Dir(exePath))

	// user vendor directory, which is part of an implied local mirror.
	arch := runtime.GOOS + "_" + runtime.GOARCH
	vendorDir := filepath.Join(rootDir, "terraform.d", "plugins", arch)
	dirs = append(dirs, vendorDir)

	// auto installed directory for Terraform v0.14+
	// The path contains a fully qualified provider name and version.
	// .terraform/providers/registry.terraform.io/hashicorp/aws/3.17.0/darwin_amd64
	// So we peek a lock file (.terraform.lock.hcl) and build a
	// path of plugin directories.
	// We don't check the `plugin_cache_dir` setting in the Terraform CLI
	// configuration or the TF_PLUGIN_CACHE_DIR environmental variable here,
	// https://github.com/hashicorp/terraform/blob/v0.14.0-rc1/website/docs/commands/cli-config.html.markdown#provider-plugin-cache
	// but even if it is configured, the lock file is stored under
	// the root directory and provider binaries are symlinked to the
	// cache directory when running terraform init.
	autoInstalledDirsV014, err := newLockFile(filepath.Join(rootDir, ".terraform.lock.hcl")).pluginDirs()
	if err != nil {
		return []string{}, err
	}

	dirs = append(dirs, autoInstalledDirsV014...)

	// auto installed directory for Terraform v0.13
	// The path contains a fully qualified provider name and version.
	// .terraform/plugins/registry.terraform.io/hashicorp/aws/2.67.0/darwin_amd64
	// So we peek a selection file (.terraform/plugins/selections.json) and build a
	// path of plugin directories.
	// We don't check the `plugin_cache_dir` setting in the Terraform CLI
	// configuration or the TF_PLUGIN_CACHE_DIR environmental variable here,
	// https://github.com/hashicorp/terraform/blob/v0.13.0-beta2/website/docs/commands/cli-config.html.markdown#provider-plugin-cache
	// but even if it is configured, the selection file is stored under
	// .terraform/plugins directory and provider binaries are symlinked to the
	// cache directory when running terraform init.
	autoInstalledDirsV013, err := newSelectionFile(filepath.Join(rootDir, ".terraform", "plugins", "selections.json")).pluginDirs()
	if err != nil {
		return []string{}, err
	}

	dirs = append(dirs, autoInstalledDirsV013...)

	// auto installed directory for Terraform < v0.13
	legacyAutoInstalledDir := filepath.Join(rootDir, ".terraform", "plugins", arch)
	dirs = append(dirs, legacyAutoInstalledDir)

	// global plugin directory
	homeDir, err := homedir.Dir()
	if err != nil {
		return []string{}, fmt.Errorf("Failed to get home dir: %s", err)
	}
	configDir := filepath.Join(homeDir, ".terraform.d", "plugins")
	dirs = append(dirs, configDir)
	dirs = append(dirs, filepath.Join(configDir, arch))

	// We don't check a provider_installation block in the Terraform CLI configuration.
	// Because it needs parse HCL and a lot of considerations to implement it
	// precisely such as include/exclude rules and a packed layout.
	// https://github.com/hashicorp/terraform/blob/v0.13.0-beta2/website/docs/commands/cli-config.html.markdown#explicit-installation-method-configuration

	// For completeness, we also should check implied local mirror directories.
	// https://github.com/hashicorp/terraform/blob/v0.13.0-beta2/website/docs/commands/cli-config.html.markdown#implied-local-mirror-directories
	// The set of directies depends on the operating system where you are running Terraform.
	// but we cannot enough test for them without test environments,
	// so we intentionally don't support it for now.
	// - Windows: %APPDATA%/HashiCorp/Terraform/plugins
	// - Mac OS X: ~/Library/Application Support/io.terraform/plugins and
	//   /Library/Application Support/io.terraform/plugins
	// - Linux and other Unix-like systems: Terraform implements the XDG Base
	//   Directory specification and appends terraform/plugins to all of the
	//   specified data directories. Without any XDG environment variables set,
	//   Terraform will use ~/.local/share/terraform/plugins,
	//   /usr/local/share/terraform/plugins, and /usr/share/terraform/plugins.

	// GOPATH
	// This is not included in the Terraform, but for convenience.
	gopath := build.Default.GOPATH
	dirs = append(dirs, filepath.Join(gopath, "bin"))

	log.Printf("[DEBUG] plugin dirs: %#v", dirs)
	return dirs, nil
}
