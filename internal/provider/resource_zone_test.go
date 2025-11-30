package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"snitchdns-tf/internal/testcontainer"
)

// TestAccZoneResource tests the Zone resource CRUD operations
func TestAccZoneResource(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	ctx := context.Background()

	// Start SnitchDNS test container
	container, err := testcontainer.NewSnitchDNSContainer(ctx, testcontainer.SnitchDNSContainerRequest{
		ExposePorts: true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	t.Logf("Container started at: %s", container.HTTPHost)
	t.Logf("API Key: %s", container.APIKey[:10]+"...")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(container),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccZoneResourceConfig(container, "test.example.com", true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snitchdns_zone.test", "domain", "test.example.com"),
					resource.TestCheckResourceAttr("snitchdns_zone.test", "active", "true"),
					resource.TestCheckResourceAttr("snitchdns_zone.test", "catch_all", "false"),
					resource.TestCheckResourceAttr("snitchdns_zone.test", "forwarding", "false"),
					resource.TestCheckResourceAttr("snitchdns_zone.test", "regex", "false"),
					resource.TestCheckResourceAttrSet("snitchdns_zone.test", "id"),
					resource.TestCheckResourceAttrSet("snitchdns_zone.test", "created_at"),
					resource.TestCheckResourceAttrSet("snitchdns_zone.test", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "snitchdns_zone.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccZoneResourceConfig(container, "test.example.com", false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snitchdns_zone.test", "domain", "test.example.com"),
					resource.TestCheckResourceAttr("snitchdns_zone.test", "active", "false"),
					resource.TestCheckResourceAttr("snitchdns_zone.test", "catch_all", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccZoneResource_WithTags tests Zone resource with tags
func TestAccZoneResource_WithTags(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	ctx := context.Background()

	container, err := testcontainer.NewSnitchDNSContainer(ctx, testcontainer.SnitchDNSContainerRequest{
		ExposePorts: true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer container.Terminate(ctx)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(container),
		Steps: []resource.TestStep{
			{
				Config: testAccZoneResourceConfigWithTags(container, "tagged.example.com", []string{"production", "web"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snitchdns_zone.test", "domain", "tagged.example.com"),
					resource.TestCheckResourceAttr("snitchdns_zone.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("snitchdns_zone.test", "tags.0", "production"),
					resource.TestCheckResourceAttr("snitchdns_zone.test", "tags.1", "web"),
				),
			},
		},
	})
}

// testAccZoneResourceConfig generates HCL configuration for testing
func testAccZoneResourceConfig(container *testcontainer.SnitchDNSContainer, domain string, active bool, catchAll bool) string {
	return fmt.Sprintf(`
provider "snitchdns" {
  api_url = %[1]q
  api_key = %[2]q
}

resource "snitchdns_zone" "test" {
  domain     = %[3]q
  active     = %[4]t
  catch_all  = %[5]t
  forwarding = false
  regex      = false
}
`, container.GetAPIEndpoint(), container.APIKey, domain, active, catchAll)
}

// testAccZoneResourceConfigWithTags generates HCL configuration with tags
func testAccZoneResourceConfigWithTags(container *testcontainer.SnitchDNSContainer, domain string, tags []string) string {
	tagsHCL := "["
	for i, tag := range tags {
		if i > 0 {
			tagsHCL += ", "
		}
		tagsHCL += fmt.Sprintf("%q", tag)
	}
	tagsHCL += "]"

	return fmt.Sprintf(`
provider "snitchdns" {
  api_url = %[1]q
  api_key = %[2]q
}

resource "snitchdns_zone" "test" {
  domain     = %[3]q
  active     = true
  catch_all  = false
  forwarding = false
  regex      = false
  tags       = %[4]s
}
`, container.GetAPIEndpoint(), container.APIKey, domain, tagsHCL)
}

// testAccCheckZoneDestroy verifies the zone has been destroyed
func testAccCheckZoneDestroy(s *terraform.State) error {
	// This will be implemented when we have the client
	return nil
}
