package azurerm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMContainerRegistryName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "four",
			ErrCount: 1,
		},
		{
			Value:    "5five",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 1,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloWorld",
			ErrCount: 0,
		},
		{
			Value:    "helloworld12",
			ErrCount: 0,
		},
		{
			Value:    "hello@world",
			ErrCount: 1,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd3324120",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd33241202",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqfjjfewsqwcdw21ddwqwd3324120",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateAzureRMContainerRegistryName(tc.Value, "azurerm_container_registry")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Container Registry Name to trigger a validation error: %v", errors)
		}
	}
}

func TestAccAzureRMContainerRegistry_basicClassic(t *testing.T) {
	resourceName := "azurerm_container_registry.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMContainerRegistry_basicUnmanaged(ri, rs, testLocation(), "Classic")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(resourceName),
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

func TestAccAzureRMContainerRegistry_basicBasic(t *testing.T) {
	resourceName := "azurerm_container_registry.test"
	ri := acctest.RandInt()
	config := testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Basic")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(resourceName),
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

func TestAccAzureRMContainerRegistry_basicStandard(t *testing.T) {
	resourceName := "azurerm_container_registry.test"
	ri := acctest.RandInt()
	config := testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Standard")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(resourceName),
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

func TestAccAzureRMContainerRegistry_basicPremium(t *testing.T) {
	resourceName := "azurerm_container_registry.test"
	ri := acctest.RandInt()
	config := testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Premium")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(resourceName),
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

func TestAccAzureRMContainerRegistry_basicBasicUpgradePremium(t *testing.T) {
	resourceName := "azurerm_container_registry.test"
	ri := acctest.RandInt()
	config := testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Basic")
	configUpdated := testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Premium")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Basic"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Premium"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_complete(t *testing.T) {
	resourceName := "azurerm_container_registry.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMContainerRegistry_complete(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(resourceName),
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

func TestAccAzureRMContainerRegistry_update(t *testing.T) {
	resourceName := "azurerm_container_registry.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccAzureRMContainerRegistry_complete(ri, rs, location)
	updatedConfig := testAccAzureRMContainerRegistry_completeUpdated(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_geoReplication(t *testing.T) {
	dataSourceName := "azurerm_container_registry.test"
	skuPremium := "Premium"
	skuBasic := "Basic"
	ri := acctest.RandInt()
	containerRegistryName := fmt.Sprintf("testacccr%d", ri)
	resourceGroupName := fmt.Sprintf("testAccRg-%d", ri)
	config := testAccAzureRMContainerRegistry_geoReplication(ri, testLocation(), skuPremium, `"eastus", "westus"`)
	updatedConfig := testAccAzureRMContainerRegistry_geoReplication(ri, testLocation(), skuPremium, `"centralus", "eastus"`)
	updatedConfigWithNoLocation := testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation(ri, testLocation(), skuPremium)
	updatedConfigBasicSku := testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation(ri, testLocation(), skuBasic)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			// first config creates an ACR with locations
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "sku", skuPremium),
					resource.TestCheckResourceAttr(dataSourceName, "georeplication_locations.#", "2"),
					testCheckAzureRMContainerRegistryExists(dataSourceName),
					testCheckAzureRMContainerRegistryGeoreplications(dataSourceName, skuPremium, []string{`"eastus"`, `"westus"`}),
				),
			},
			// second config udpates the ACR with updated locations
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "sku", skuPremium),
					resource.TestCheckResourceAttr(dataSourceName, "georeplication_locations.#", "2"),
					testCheckAzureRMContainerRegistryExists(dataSourceName),
					testCheckAzureRMContainerRegistryGeoreplications(dataSourceName, skuPremium, []string{`"eastus"`, `"centralus"`}),
				),
			},
			// third config udpates the ACR with no location
			{
				Config: updatedConfigWithNoLocation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "sku", skuPremium),
					testCheckAzureRMContainerRegistryExists(dataSourceName),
					testCheckAzureRMContainerRegistryGeoreplications(dataSourceName, skuPremium, nil),
				),
			},
			// fourth config updates an ACR with replicas
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "sku", skuPremium),
					resource.TestCheckResourceAttr(dataSourceName, "georeplication_locations.#", "2"),
					testCheckAzureRMContainerRegistryExists(dataSourceName),
					testCheckAzureRMContainerRegistryGeoreplications(dataSourceName, skuPremium, []string{`"eastus"`, `"westus"`}),
				),
			},
			// fifth config updates the SKU to basic and no replicas (should remove the existing replicas if any)
			{
				Config: updatedConfigBasicSku,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "sku", skuBasic),
					testCheckAzureRMContainerRegistryExists(dataSourceName),
					testCheckAzureRMContainerRegistryGeoreplications(dataSourceName, skuBasic, nil),
				),
			},
		},
	})
}

func testCheckAzureRMContainerRegistryDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).containerRegistryClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_registry" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testCheckAzureRMContainerRegistryExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Container Registry: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).containerRegistryClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on containerRegistryClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Container Registry %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMContainerRegistryGeoreplications(registryName string, sku string, expectedLocations []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[registryName]
		if !ok {
			return fmt.Errorf("Not found: %s", registryName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Container Registry: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).containerRegistryReplicationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.List(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on containerRegistryClient: %+v", err)
		}

		georeplicationValues := resp.Values()
		expectedLocationsCount := len(expectedLocations) + 1 // the main location is returned by the API as a geolocation for replication.

		// if Sku is not premium, listing the geo-replications locations returns an empty list
		if strings.ToLower(sku) != strings.ToLower(string(containerregistry.Premium)) {
			expectedLocationsCount = 0
		}

		actualLocationsCount := len(georeplicationValues)

		if expectedLocationsCount != actualLocationsCount {
			return fmt.Errorf("Bad: Container Registry %q (resource group: %q) expected locations count is %d, actual location count is %d", name, resourceGroup, expectedLocationsCount, actualLocationsCount)
		}

		return nil
	}
}

func testAccAzureRMContainerRegistry_basicManaged(rInt int, location string, sku string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "%s"
}
`, rInt, location, rInt, sku)
}

func testAccAzureRMContainerRegistry_basicUnmanaged(rInt int, rStr string, location string, sku string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "%s"
  storage_account_id  = "${azurerm_storage_account.test.id}"
}
`, rInt, location, rStr, rInt, sku)
}

func testAccAzureRMContainerRegistry_complete(rInt int, rStr string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  admin_enabled       = false
  sku                 = "Classic"
  storage_account_id  = "${azurerm_storage_account.test.id}"

  tags {
    environment = "production"
  }
}
`, rInt, location, rStr, rInt)
}

func testAccAzureRMContainerRegistry_completeUpdated(rInt int, rStr string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  admin_enabled       = true
  sku                 = "Classic"
  storage_account_id  = "${azurerm_storage_account.test.id}"

  tags {
    environment = "production"
  }
}
`, rInt, location, rStr, rInt)
}

func testAccAzureRMContainerRegistry_geoReplication(rInt int, location string, sku string, georeplicationLocations string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testAccRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                   = "testacccr%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  location               = "${azurerm_resource_group.test.location}"
  sku                    = "%s"
  georeplication_locations = [%s]
}
`, rInt, location, rInt, sku, georeplicationLocations)
}

func testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation(rInt int, location string, sku string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testAccRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                   = "testacccr%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  location               = "${azurerm_resource_group.test.location}"
  sku                    = "%s"
}
`, rInt, location, rInt, sku)
}
