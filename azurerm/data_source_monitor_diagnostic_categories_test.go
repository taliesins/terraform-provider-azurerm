package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceArmMonitorDiagnosticCategories_appService(t *testing.T) {
	dataSourceName := "data.azurerm_monitor_diagnostic_categories.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmMonitorDiagnosticCategories_appService(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "metrics.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "logs.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceArmMonitorDiagnosticCategories_storageAccount(t *testing.T) {
	dataSourceName := "data.azurerm_monitor_diagnostic_categories.test"
	rs := acctest.RandString(8)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmMonitorDiagnosticCategories_storageAccount(rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "metrics.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "logs.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceArmMonitorDiagnosticCategories_appService(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

data "azurerm_monitor_diagnostic_categories" "test" {
  resource_id = "${azurerm_app_service.test.id}"
}
`, rInt, location, rInt, rInt)
}

func testAccDataSourceArmMonitorDiagnosticCategories_storageAccount(rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%s"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags {
    environment = "staging"
  }
}

data "azurerm_monitor_diagnostic_categories" "test" {
  resource_id = "${azurerm_storage_account.test.id}"
}
`, rString, location, rString)
}
