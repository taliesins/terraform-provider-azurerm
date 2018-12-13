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

func TestAccAzureRMEventHubNamespaceAuthorizationRule_listen(t *testing.T) {
	testAccAzureRMEventHubNamespaceAuthorizationRule(t, true, false, false)
}

func TestAccAzureRMEventHubNamespaceAuthorizationRule_send(t *testing.T) {
	testAccAzureRMEventHubNamespaceAuthorizationRule(t, false, true, false)
}

func TestAccAzureRMEventHubNamespaceAuthorizationRule_listensend(t *testing.T) {
	testAccAzureRMEventHubNamespaceAuthorizationRule(t, true, true, false)
}

func TestAccAzureRMEventHubNamespaceAuthorizationRule_manage(t *testing.T) {
	testAccAzureRMEventHubNamespaceAuthorizationRule(t, true, true, true)
}

func testAccAzureRMEventHubNamespaceAuthorizationRule(t *testing.T, listen, send, manage bool) {
	resourceName := "azurerm_eventhub_namespace_authorization_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceAuthorizationRule_base(acctest.RandInt(), testLocation(), listen, send, manage),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_name"),
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

func TestAccAzureRMEventHubNamespaceAuthorizationRule_rightsUpdate(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace_authorization_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceAuthorizationRule_base(acctest.RandInt(), testLocation(), true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttr(resourceName, "send", "false"),
					resource.TestCheckResourceAttr(resourceName, "manage", "false"),
				),
			},
			{
				Config: testAccAzureRMEventHubNamespaceAuthorizationRule_base(acctest.RandInt(), testLocation(), true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceAuthorizationRuleExists(resourceName),
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

func testCheckAzureRMEventHubNamespaceAuthorizationRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).eventHubNamespacesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_authorization_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMEventHubNamespaceAuthorizationRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Event Hub: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).eventHubNamespacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Event Hub Namespace Authorization Rule %q (namespace %q / resource group: %q) does not exist", name, namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on eventHubClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMEventHubNamespaceAuthorizationRule_base(rInt int, location string, listen, send, manage bool) string {
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

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  listen              = %[3]t
  send                = %[4]t
  manage              = %[5]t
}
`, rInt, location, listen, send, manage)
}
