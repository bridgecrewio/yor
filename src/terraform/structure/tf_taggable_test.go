package structure

import (
	"testing"
)

func TestTaggableResourceShouldBeTaggable(t *testing.T) {
	expectedResources := []string{
		"aws_vpc",
		"azurerm_resource_group",
		"google_compute_instance",
	}
	resources := make(map[string]struct{}, 0)
	for _, r := range TfTaggableResourceTypes {
		resources[r] = struct{}{}
	}
	for _, r := range expectedResources {
		t.Run(r, func(t *testing.T) {
			_, ok := resources[r]
			if !ok {
				t.Errorf("%s should be taggable", r)
			}
		})
	}
}

func TestNonTaggableResourceShouldNotBeTaggable(t *testing.T) {
	nonTaggableResources := []string{
		"aws_vpc_endpoint_subnet_association",
		"azurerm_virtual_machine_data_disk_attachment",
		"google_compute_attached_disk",
	}
	resources := make(map[string]struct{}, 0)
	for _, r := range TfTaggableResourceTypes {
		resources[r] = struct{}{}
	}
	for _, r := range nonTaggableResources {
		t.Run(r, func(t *testing.T) {
			_, ok := resources[r]
			if ok {
				t.Errorf("%s should not be taggable", r)
			}
		})
	}
}

func TestResourceThatDeprecatedTagsInLatestProviderShouldNotBeTaggable(t *testing.T) {
	resource := "azurerm_log_analytics_linked_service"
	for _, n := range TfTaggableResourceTypes {
		if resource == n {
			t.Errorf("`azurerm_log_analytics_linked_service`'s tags has been removed since AzureRM 3.0")
		}
	}
}

func TestTagsWithIncorrectTypeShouldNotBeTaggable(t *testing.T) {
	resource := "azurerm_api_management_property"
	for _, r := range TfTaggableResourceTypes {
		if r == resource {
			t.Errorf("`azurerm_api_management_property`'s tags is a list so it should not be taggable")
		}
	}
}

func TestDeprecatedTaggableResourceShouldBeTaggable(t *testing.T) {
	deprecatedResource := "azurerm_data_lake_analytics_account"
	taggable := false
	for _, r := range TfTaggableResourceTypes {
		if r == deprecatedResource {
			taggable = true
		}
	}
	if !taggable {
		t.Errorf("`azurerm_data_lake_analytics_account` should be taggable")
	}
}

func TestKnownUntaggableResourcesShouldNotBeTaggable(t *testing.T) {
	taggable := make(map[string]struct{})
	for _, n := range TfTaggableResourceTypes {
		taggable[n] = struct{}{}
	}
	for _, r := range unsupportedTerraformBlocks {
		_, ok := taggable[r]
		if ok {
			t.Errorf("%s should not be taggable", r)
		}
	}
}

func BenchmarkInit(b *testing.B) {
	loadSchema()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loadSchema()
	}
}
