package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAppServiceActiveSlot_basic(t *testing.T) {
	resourceName := "azurerm_app_service_active_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceActiveSlot_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// Destroy actually does nothing so we just return nil
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app_service_slot_name", fmt.Sprintf("acctestASSlot-%d", ri)),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceActiveSlot_update(t *testing.T) {
	resourceName := "azurerm_app_service_active_slot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceActiveSlot_update(ri, testLocation())
	config2 := testAccAzureRMAppServiceActiveSlot_updated(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// Destroy actually does nothing so we just return nil
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app_service_slot_name", fmt.Sprintf("acctestASSlot-%d", ri)),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app_service_slot_name", fmt.Sprintf("acctestASSlot2-%d", ri)),
				),
			},
		},
	})
}

func testAccAzureRMAppServiceActiveSlot_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  app_service_name    = "${azurerm_app_service.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_active_slot" "test" {
  resource_group_name   = "${azurerm_resource_group.test.name}"
  app_service_name      = "${azurerm_app_service.test.name}"
  app_service_slot_name = "${azurerm_app_service_slot.test.name}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceActiveSlot_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  app_service_name    = "${azurerm_app_service.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test2" {
  name                = "acctestASSlot2-%d"
  app_service_name    = "${azurerm_app_service.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_active_slot" "test" {
  resource_group_name   = "${azurerm_resource_group.test.name}"
  app_service_name      = "${azurerm_app_service.test.name}"
  app_service_slot_name = "${azurerm_app_service_slot.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceActiveSlot_updated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  app_service_name    = "${azurerm_app_service.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_slot" "test2" {
  name                = "acctestASSlot2-%d"
  app_service_name    = "${azurerm_app_service.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_active_slot" "test" {
  resource_group_name   = "${azurerm_resource_group.test.name}"
  app_service_name      = "${azurerm_app_service.test.name}"
  app_service_slot_name = "${azurerm_app_service_slot.test2.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
