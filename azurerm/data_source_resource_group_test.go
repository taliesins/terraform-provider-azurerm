package azurerm

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMResourceGroup_basic(t *testing.T) {
	ri := acctest.RandInt()
	name := fmt.Sprintf("acctestRg_%d", ri)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMResourceGroupBasic(name, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.azurerm_resource_group.test", "name", name),
					resource.TestCheckResourceAttr("data.azurerm_resource_group.test", "location", azureRMNormalizeLocation(location)),
					resource.TestCheckResourceAttr("data.azurerm_resource_group.test", "tags.%", "1"),
					resource.TestCheckResourceAttr("data.azurerm_resource_group.test", "tags.env", "test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMResourceGroupBasic(name string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"

  tags {
    env = "test"
  }
}

data "azurerm_resource_group" "test" {
  name = "${azurerm_resource_group.test.name}"
}
`, name, location)
}
