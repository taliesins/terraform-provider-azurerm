package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMDevTestLab_basic(t *testing.T) {
	dataSourceName := "data.azurerm_dev_test_lab.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDevTestLab_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "storage_type", "Premium"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMDevTestLab_complete(t *testing.T) {
	dataSourceName := "data.azurerm_dev_test_lab.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDevTestLab_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "storage_type", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.Hello", "World"),
				),
			},
		},
	})
}

func testAccDataSourceDevTestLab_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

data "azurerm_dev_test_lab" "test" {
  name                = "${azurerm_dev_test_lab.test.name}"
  resource_group_name = "${azurerm_dev_test_lab.test.resource_group_name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceDevTestLab_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  storage_type        = "Standard"

  tags {
    "Hello" = "World"
  }
}

data "azurerm_dev_test_lab" "test" {
  name                = "${azurerm_dev_test_lab.test.name}"
  resource_group_name = "${azurerm_dev_test_lab.test.resource_group_name}"
}
`, rInt, location, rInt)
}
