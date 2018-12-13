package azurerm

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDataLakeStore_basic(t *testing.T) {
	resourceName := "azurerm_data_lake_store.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tier", "Consumption"),
					resource.TestCheckResourceAttr(resourceName, "encryption_state", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "encryption_type", "ServiceManaged"),
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

func TestAccAzureRMDataLakeStore_tier(t *testing.T) {
	resourceName := "azurerm_data_lake_store.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_tier(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tier", "Commitment_1TB"),
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

func TestAccAzureRMDataLakeStore_encryptionDisabled(t *testing.T) {
	resourceName := "azurerm_data_lake_store.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_encryptionDisabled(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "encryption_state", "Disabled"),
					resource.TestCheckResourceAttr(resourceName, "encryption_type", ""),
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

func TestAccAzureRMDataLakeStore_firewallUpdate(t *testing.T) {
	resourceName := "azurerm_data_lake_store.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_firewall(ri, location, "Enabled", "Enabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firewall_state", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "firewall_allow_azure_ips", "Enabled"),
				),
			},
			{
				Config: testAccAzureRMDataLakeStore_firewall(ri, location, "Enabled", "Disabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firewall_state", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "firewall_allow_azure_ips", "Disabled"),
				),
			},
			{
				Config: testAccAzureRMDataLakeStore_firewall(ri, location, "Disabled", "Enabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firewall_state", "Disabled"),
					resource.TestCheckResourceAttr(resourceName, "firewall_allow_azure_ips", "Enabled"),
				),
			},
			{
				Config: testAccAzureRMDataLakeStore_firewall(ri, location, "Disabled", "Disabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firewall_state", "Disabled"),
					resource.TestCheckResourceAttr(resourceName, "firewall_allow_azure_ips", "Disabled"),
				),
			},
		},
	})
}

func TestAccAzureRMDataLakeStore_withTags(t *testing.T) {
	resourceName := "azurerm_data_lake_store.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_withTags(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDataLakeStore_withTagsUpdate(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
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

func testCheckAzureRMDataLakeStoreExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		accountName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for data lake store: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).dataLakeStoreAccountClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Bad: Get on dataLakeStoreAccountClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Date Lake Store %q (resource group: %q) does not exist", accountName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataLakeStoreDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dataLakeStoreAccountClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_lake_store" {
			continue
		}

		accountName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, accountName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Data Lake Store still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDataLakeStore_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, rInt, location, strconv.Itoa(rInt)[0:15])
}

func testAccAzureRMDataLakeStore_tier(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  tier                = "Commitment_1TB"
}
`, rInt, location, strconv.Itoa(rInt)[0:15])
}

func testAccAzureRMDataLakeStore_encryptionDisabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  encryption_state    = "Disabled"
}
`, rInt, location, strconv.Itoa(rInt)[0:15])
}

func testAccAzureRMDataLakeStore_firewall(rInt int, location string, firewallState string, firewallAllowAzureIPs string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                     = "acctest%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  firewall_state           = "%s"
  firewall_allow_azure_ips = "%s"
}
`, rInt, location, strconv.Itoa(rInt)[0:15], firewallState, firewallAllowAzureIPs)
}

func testAccAzureRMDataLakeStore_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  tags {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, strconv.Itoa(rInt)[0:15])
}

func testAccAzureRMDataLakeStore_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  tags {
    environment = "staging"
  }
}
`, rInt, location, strconv.Itoa(rInt)[0:15])
}
