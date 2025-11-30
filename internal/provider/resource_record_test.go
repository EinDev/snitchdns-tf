package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"snitchdns-tf/internal/testcontainer"
)

// TestAccRecordResource tests the Record resource CRUD operations
func TestAccRecordResource(t *testing.T) {
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
			// Create zone first, then create A record
			{
				Config: testAccRecordResourceConfigA(container, "record-test.example.com", "192.168.1.1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("snitchdns_record.test", "id"),
					resource.TestCheckResourceAttrSet("snitchdns_record.test", "zone_id"),
					resource.TestCheckResourceAttr("snitchdns_record.test", "type", "A"),
					resource.TestCheckResourceAttr("snitchdns_record.test", "cls", "IN"),
					resource.TestCheckResourceAttr("snitchdns_record.test", "ttl", "300"),
					resource.TestCheckResourceAttr("snitchdns_record.test", "active", "true"),
					resource.TestCheckResourceAttr("snitchdns_record.test", "data.address", "192.168.1.1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "snitchdns_record.test",
				ImportState:       true,
				ImportStateIdFunc: testAccRecordImportStateIdFunc,
				ImportStateVerify: true,
			},
			// Update the A record IP address
			{
				Config: testAccRecordResourceConfigA(container, "record-test.example.com", "192.168.1.2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snitchdns_record.test", "data.address", "192.168.1.2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccRecordResource_CNAME tests CNAME record creation
func TestAccRecordResource_CNAME(t *testing.T) {
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
				Config: testAccRecordResourceConfigCNAME(container, "cname-test.example.com", "target.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snitchdns_record.test", "type", "CNAME"),
					resource.TestCheckResourceAttr("snitchdns_record.test", "data.name", "target.example.com"),
				),
			},
		},
	})
}

// testAccRecordImportStateIdFunc returns the import ID in format "zone_id:record_id"
func testAccRecordImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["snitchdns_record.test"]
	if !ok {
		return "", fmt.Errorf("Resource not found")
	}

	zoneID := rs.Primary.Attributes["zone_id"]
	recordID := rs.Primary.ID

	return fmt.Sprintf("%s:%s", zoneID, recordID), nil
}

// testAccRecordResourceConfigA generates HCL configuration for A record testing
func testAccRecordResourceConfigA(container *testcontainer.SnitchDNSContainer, domain string, address string) string {
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
}

resource "snitchdns_record" "test" {
  zone_id = snitchdns_zone.test.id
  type    = "A"
  cls     = "IN"
  ttl     = 300
  active  = true

  data = {
    address = %[4]q
  }
}
`, container.GetAPIEndpoint(), container.APIKey, domain, address)
}

// testAccRecordResourceConfigCNAME generates HCL configuration for CNAME record testing
func testAccRecordResourceConfigCNAME(container *testcontainer.SnitchDNSContainer, domain string, target string) string {
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
}

resource "snitchdns_record" "test" {
  zone_id = snitchdns_zone.test.id
  type    = "CNAME"
  cls     = "IN"
  ttl     = 300
  active  = true

  data = {
    name = %[4]q
  }
}
`, container.GetAPIEndpoint(), container.APIKey, domain, target)
}
