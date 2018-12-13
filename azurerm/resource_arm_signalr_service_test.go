package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMSignalRService_basic(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSignalRService_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSignalRService_standard(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSignalRService_standardWithCapacity(ri, testLocation(), 1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSignalRService_standardWithCap2(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSignalRService_standardWithCapacity(ri, testLocation(), 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSignalRService_skuUpdate(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := acctest.RandInt()
	location := testLocation()
	freeConfig := testAccAzureRMSignalRService_basic(ri, location)
	standardConfig := testAccAzureRMSignalRService_standardWithCapacity(ri, location, 1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: freeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
			{
				Config: standardConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
			{
				Config: freeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
		},
	})
}

func TestAccAzureRMSignalRService_capacityUpdate(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := acctest.RandInt()
	location := testLocation()
	standardConfig := testAccAzureRMSignalRService_standardWithCapacity(ri, location, 1)
	standardCap5Config := testAccAzureRMSignalRService_standardWithCapacity(ri, location, 5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: standardConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
			{
				Config: standardCap5Config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
			{
				Config: standardConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
		},
	})
}

func TestAccAzureRMSignalRService_skuAndCapacityUpdate(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := acctest.RandInt()
	location := testLocation()
	freeConfig := testAccAzureRMSignalRService_basic(ri, location)
	standardConfig := testAccAzureRMSignalRService_standardWithCapacity(ri, location, 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: freeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
			{
				Config: standardConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
			{
				Config: freeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
				),
			},
		},
	})
}

func testAccAzureRMSignalRService_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku {
    name     = "Free_F1"
    capacity = 1
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMSignalRService_standardWithCapacity(rInt int, location string, capacity int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku {
    name     = "Standard_S1"
    capacity = %d
  }
}
`, rInt, location, rInt, capacity)
}

func testCheckAzureRMSignalRServiceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).signalRClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_signalr_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("SignalR service still exists:\n%#v", resp)
		}
	}
	return nil
}

func testCheckAzureRMSignalRServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for SignalR service: %s", resourceName)
		}

		conn := testAccProvider.Meta().(*ArmClient).signalRClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, resourceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on signalRClient: %+v", err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: SignalR service %q (resource group: %q) does not exist", resourceName, resourceGroup)
		}

		return nil
	}
}
