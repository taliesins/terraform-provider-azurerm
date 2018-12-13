package azurerm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventHubAuthorizationRule_listen(t *testing.T) {
	testAccAzureRMEventHubAuthorizationRule(t, true, false, false)
}

func TestAccAzureRMEventHubAuthorizationRule_send(t *testing.T) {
	testAccAzureRMEventHubAuthorizationRule(t, false, true, false)
}

func TestAccAzureRMEventHubAuthorizationRule_listensend(t *testing.T) {
	testAccAzureRMEventHubAuthorizationRule(t, true, true, false)
}

func TestAccAzureRMEventHubAuthorizationRule_manage(t *testing.T) {
	testAccAzureRMEventHubAuthorizationRule(t, true, true, true)
}

func testAccAzureRMEventHubAuthorizationRule(t *testing.T, listen, send, manage bool) {
	resourceName := "azurerm_eventhub_authorization_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubAuthorizationRule_base(acctest.RandInt(), testLocation(), listen, send, manage),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub_name"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(resourceName, "listen", strconv.FormatBool(listen)),
					resource.TestCheckResourceAttr(resourceName, "send", strconv.FormatBool(send)),
					resource.TestCheckResourceAttr(resourceName, "manage", strconv.FormatBool(manage)),
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

func TestAccAzureRMEventHubAuthorizationRule_rightsUpdate(t *testing.T) {
	resourceName := "azurerm_eventhub_authorization_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubAuthorizationRule_base(acctest.RandInt(), testLocation(), true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttr(resourceName, "send", "false"),
					resource.TestCheckResourceAttr(resourceName, "manage", "false"),
				),
			},
			{
				Config: testAccAzureRMEventHubAuthorizationRule_base(acctest.RandInt(), testLocation(), true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttr(resourceName, "send", "true"),
					resource.TestCheckResourceAttr(resourceName, "manage", "true"),
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

func testCheckAzureRMEventHubAuthorizationRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).eventHubClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_authorization_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		eventHubName := rs.Primary.Attributes["eventhub_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, eventHubName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMEventHubAuthorizationRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		eventHubName := rs.Primary.Attributes["eventhub_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Event Hub: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).eventHubClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, eventHubName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Event Hub Authorization Rule %q (eventhub %s / namespace %s / resource group: %s) does not exist", name, eventHubName, namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on eventHubClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMEventHubAuthorizationRule_base(rInt int, location string, listen, send, manage bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[1]d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  eventhub_name       = "${azurerm_eventhub.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  listen              = %[3]t
  send                = %[4]t
  manage              = %[5]t
}
`, rInt, location, listen, send, manage)
}
