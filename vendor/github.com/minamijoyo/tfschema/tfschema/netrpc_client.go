package tfschema

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sort"

	"github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	tfplugin "github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/hashicorp/terraform/terraform"
)

// NetRPCClient implements Client interface.
// This implementaion is for Terraform v0.11.
type NetRPCClient struct {
	// provider is a resource provider of Terraform.
	provider terraform.ResourceProvider
	// pluginClient is a pointer to plugin client instance.
	// The type of pluginClient is
	// *github.com/hashicorp/terraform/vendor/github.com/hashicorp/go-plugin.Client.
	// But, we cannot import the vendor version of go-plugin using terraform.
	// So, we store this as interface{}, and use it by reflection.
	pluginClient interface{}
}

// NewNetRPCClient creates a new NetRPCClient instance.
func NewNetRPCClient(providerName string, options Option) (Client, error) {
	// find a provider plugin
	pluginMeta, err := findPlugin("provider", providerName, options.RootDir)
	if err != nil {
		return nil, err
	}

	// create a plugin client config
	config := newNetRPCClientConfig(pluginMeta)
	if options.Logger != nil {
		config.Logger = options.Logger
	}

	// initialize a plugin client.
	pluginClient := plugin.NewClient(config)
	client, err := pluginClient.Client()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize NetRPC plugin: %s", err)
	}

	// create a new ResourceProvider.
	raw, err := client.Dispense(tfplugin.ProviderPluginName)
	if err != nil {
		return nil, fmt.Errorf("Failed to dispense NetRPC plugin: %s", err)
	}

	switch provider := raw.(type) {
	// For Terraform v0.11
	case *tfplugin.ResourceProvider:
		return &NetRPCClient{
			provider:     provider,
			pluginClient: pluginClient,
		}, nil

	default:
		return nil, fmt.Errorf("Failed to type cast NetRPC plugin r: %+v", raw)
	}
}

// newNetRPCClientConfig returns a default plugin client config for Terraform v0.11.
func newNetRPCClientConfig(pluginMeta *discovery.PluginMeta) *plugin.ClientConfig {
	// Note that we depends on Terraform v0.12 library
	// and cannot simply refer the v0.11 default config.
	// So, we need to reproduce the v0.11 config manually.
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Level:  hclog.Trace,
		Output: os.Stderr,
	})

	pluginMap := map[string]plugin.Plugin{
		"provider":    &tfplugin.ResourceProviderPlugin{},
		"provisioner": &tfplugin.ResourceProvisionerPlugin{},
	}

	return &plugin.ClientConfig{
		Cmd:             exec.Command(pluginMeta.Path),
		HandshakeConfig: tfplugin.Handshake,
		Managed:         true,
		Plugins:         pluginMap,
		Logger:          logger,
	}
}

// GetProviderSchema returns a type definiton of provider schema.
func (c *NetRPCClient) GetProviderSchema() (*Block, error) {
	req := &terraform.ProviderSchemaRequest{
		ResourceTypes: []string{},
		DataSources:   []string{},
	}

	res, err := c.provider.GetSchema(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get schema from provider: %s", err)
	}

	b := NewBlock(res.Provider)
	return b, nil
}

// GetResourceTypeSchema returns a type definiton of resource type.
func (c *NetRPCClient) GetResourceTypeSchema(resourceType string) (*Block, error) {
	req := &terraform.ProviderSchemaRequest{
		ResourceTypes: []string{resourceType},
		DataSources:   []string{},
	}

	res, err := c.provider.GetSchema(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get schema from provider: %s", err)
	}

	if res.ResourceTypes[resourceType] == nil {
		return nil, fmt.Errorf("Failed to find resource type: %s", resourceType)
	}

	b := NewBlock(res.ResourceTypes[resourceType])
	return b, nil
}

// GetDataSourceSchema returns a type definiton of data source.
func (c *NetRPCClient) GetDataSourceSchema(dataSource string) (*Block, error) {
	req := &terraform.ProviderSchemaRequest{
		ResourceTypes: []string{},
		DataSources:   []string{dataSource},
	}

	res, err := c.provider.GetSchema(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get schema from provider: %s", err)
	}

	if res.DataSources[dataSource] == nil {
		return nil, fmt.Errorf("Failed to find data source: %s", dataSource)
	}

	b := NewBlock(res.DataSources[dataSource])
	return b, nil
}

// ResourceTypes returns a list of resource types.
func (c *NetRPCClient) ResourceTypes() ([]string, error) {
	res := c.provider.Resources()

	keys := make([]string, 0, len(res))
	for _, r := range res {
		keys = append(keys, r.Name)
	}

	sort.Strings(keys)
	return keys, nil
}

// DataSources returns a list of data sources.
func (c *NetRPCClient) DataSources() ([]string, error) {
	res := c.provider.DataSources()

	keys := make([]string, 0, len(res))
	for _, r := range res {
		keys = append(keys, r.Name)
	}

	sort.Strings(keys)
	return keys, nil
}

// Close kills a process of the plugin.
func (c *NetRPCClient) Close() {
	// We cannot import the vendor version of go-plugin using terraform.
	// So, we call (*go-plugin.Client).Kill() by reflection here.
	v := reflect.ValueOf(c.pluginClient).MethodByName("Kill")
	v.Call([]reflect.Value{})
}
