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

func TestAccAzureRMDataLakeAnalyticsAccount_basic(t *testing.T) {
	resourceName := "azurerm_data_lake_analytics_account.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsAccount_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tier", "Consumption"),
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

func TestAccAzureRMDataLakeAnalyticsAccount_tier(t *testing.T) {
	resourceName := "azurerm_data_lake_analytics_account.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsAccount_tier(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tier", "Commitment_100AUHours"),
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

func TestAccAzureRMDataLakeAnalyticsAccount_withTags(t *testing.T) {
	resourceName := "azurerm_data_lake_analytics_account.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsAccount_withTags(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDataLakeAnalyticsAccount_withTagsUpdate(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsAccountExists(resourceName),
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

func testCheckAzureRMDataLakeAnalyticsAccountExists(name string) resource.TestCheckFunc {
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

		conn := testAccProvider.Meta().(*ArmClient).dataLakeAnalyticsAccountClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Bad: Get on dataLakeAnalyticsAccountClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Date Lake Analytics Account %q (resource group: %q) does not exist", accountName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataLakeAnalyticsAccountDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dataLakeAnalyticsAccountClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_lake_analytics_account" {
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

		return fmt.Errorf("Data Lake Analytics Account still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDataLakeAnalyticsAccount_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  default_store_account_name = "${azurerm_data_lake_store.test.name}"
}
`, testAccAzureRMDataLakeStore_basic(rInt, location), strconv.Itoa(rInt)[0:15])
}

func testAccAzureRMDataLakeAnalyticsAccount_tier(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  tier                = "Commitment_100AUHours"

  default_store_account_name = "${azurerm_data_lake_store.test.name}"
}
`, testAccAzureRMDataLakeStore_basic(rInt, location), strconv.Itoa(rInt)[0:15])
}

func testAccAzureRMDataLakeAnalyticsAccount_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  
  default_store_account_name = "${azurerm_data_lake_store.test.name}"

  tags {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, testAccAzureRMDataLakeStore_basic(rInt, location), strconv.Itoa(rInt)[0:15])
}

func testAccAzureRMDataLakeAnalyticsAccount_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  
  default_store_account_name = "${azurerm_data_lake_store.test.name}"

  tags {
    environment = "staging"
  }
}
`, testAccAzureRMDataLakeStore_basic(rInt, location), strconv.Itoa(rInt)[0:15])
}
