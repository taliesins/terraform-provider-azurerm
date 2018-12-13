package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMySqlVirtualNetworkRule_basic(t *testing.T) {
	resourceName := "azurerm_mysql_virtual_network_rule.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySqlVirtualNetworkRule_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySqlVirtualNetworkRuleExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMMySqlVirtualNetworkRule_switchSubnets(t *testing.T) {
	resourceName := "azurerm_mysql_virtual_network_rule.test"
	ri := acctest.RandInt()

	preConfig := testAccAzureRMMySqlVirtualNetworkRule_subnetSwitchPre(ri, testLocation())
	postConfig := testAccAzureRMMySqlVirtualNetworkRule_subnetSwitchPost(ri, testLocation())

	// Create regex strings that will ensure that one subnet name exists, but not the other
	preConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet1%d)$|(subnet[^2]%d)$", ri, ri))  //subnet 1 but not 2
	postConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet2%d)$|(subnet[^1]%d)$", ri, ri)) //subnet 2 but not 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySqlVirtualNetworkRuleExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "subnet_id", preConfigRegex),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySqlVirtualNetworkRuleExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "subnet_id", postConfigRegex),
				),
			},
		},
	})
}

func TestAccAzureRMMySqlVirtualNetworkRule_disappears(t *testing.T) {
	resourceName := "azurerm_mysql_virtual_network_rule.test"
	ri := acctest.RandInt()
	config := testAccAzureRMMySqlVirtualNetworkRule_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySqlVirtualNetworkRuleExists(resourceName),
					testCheckAzureRMMySqlVirtualNetworkRuleDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMMySqlVirtualNetworkRule_multipleSubnets(t *testing.T) {
	resourceName1 := "azurerm_mysql_virtual_network_rule.rule1"
	resourceName2 := "azurerm_mysql_virtual_network_rule.rule2"
	resourceName3 := "azurerm_mysql_virtual_network_rule.rule3"
	ri := acctest.RandInt()
	config := testAccAzureRMMySqlVirtualNetworkRule_multipleSubnets(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySqlVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySqlVirtualNetworkRuleExists(resourceName1),
					testCheckAzureRMMySqlVirtualNetworkRuleExists(resourceName2),
					testCheckAzureRMMySqlVirtualNetworkRuleExists(resourceName3),
				),
			},
		},
	})
}

func testCheckAzureRMMySqlVirtualNetworkRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).mysqlVirtualNetworkRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySql Virtual Network Rule %q (Server %q / Resource Group %q) was not found", ruleName, serverName, resourceGroup)
			}

			return err
		}

		return nil
	}
}

func testCheckAzureRMMySqlVirtualNetworkRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mysql_virtual_network_rule" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).mysqlVirtualNetworkRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Bad: MySql Firewall Rule %q (Server %q / Resource Group %q) still exists: %+v", ruleName, serverName, resourceGroup, resp)
	}

	return nil
}

func testCheckAzureRMMySqlVirtualNetworkRuleDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).mysqlVirtualNetworkRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			//If the error is that the resource we want to delete does not exist in the first
			//place (404), then just return with no error.
			if response.WasNotFound(future.Response()) {
				return nil
			}

			return fmt.Errorf("Error deleting MySql Virtual Network Rule: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			//Same deal as before. Just in case.
			if response.WasNotFound(future.Response()) {
				return nil
			}

			return fmt.Errorf("Error deleting MySql Virtual Network Rule: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMMySqlVirtualNetworkRule_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mysql_server" "test" {
  name                         = "acctestmysqlsvr-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.6"
  ssl_enforcement              = "Enabled"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_mysql_virtual_network_rule" "test" {
  name                = "acctestmysqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  subnet_id           = "${azurerm_subnet.test.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMMySqlVirtualNetworkRule_subnetSwitchPre(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "subnet2%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mysql_server" "test" {
  name                         = "acctestmysqlsvr-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.6"
  ssl_enforcement              = "Enabled"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_mysql_virtual_network_rule" "test" {
  name                = "acctestmysqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  subnet_id           = "${azurerm_subnet.test1.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMMySqlVirtualNetworkRule_subnetSwitchPost(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "subnet2%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mysql_server" "test" {
  name                         = "acctestmysqlsvr-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.6"
  ssl_enforcement              = "Enabled"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_mysql_virtual_network_rule" "test" {
  name                = "acctestmysqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  subnet_id           = "${azurerm_subnet.test2.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMMySqlVirtualNetworkRule_multipleSubnets(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "vnet1" {
  name                = "acctestvnet1%d"
  address_space       = ["10.7.29.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_network" "vnet2" {
  name                = "acctestvnet2%d"
  address_space       = ["10.1.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "vnet1_subnet1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet1.name}"
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "vnet1_subnet2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet1.name}"
  address_prefix       = "10.7.29.128/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "vnet2_subnet1" {
  name                 = "acctestsubnet3%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet2.name}"
  address_prefix       = "10.1.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mysql_server" "test" {
  name                         = "acctestmysqlsvr-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.6"
  ssl_enforcement              = "Enabled"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_mysql_virtual_network_rule" "rule1" {
  name                = "acctestmysqlvnetrule1%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  subnet_id           = "${azurerm_subnet.vnet1_subnet1.id}"
}

resource "azurerm_mysql_virtual_network_rule" "rule2" {
  name                = "acctestmysqlvnetrule2%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  subnet_id           = "${azurerm_subnet.vnet1_subnet2.id}"
}

resource "azurerm_mysql_virtual_network_rule" "rule3" {
  name                = "acctestmysqlvnetrule3%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  subnet_id           = "${azurerm_subnet.vnet2_subnet1.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}
