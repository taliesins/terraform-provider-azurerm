package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultAccessPolicy_basic(t *testing.T) {
	resourceName := "azurerm_key_vault_access_policy.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultAccessPolicy_basic(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "set"),
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

func TestAccAzureRMKeyVaultAccessPolicy_multiple(t *testing.T) {
	resourceName1 := "azurerm_key_vault_access_policy.test_with_application_id"
	resourceName2 := "azurerm_key_vault_access_policy.test_no_application_id"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultAccessPolicy_multiple(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName1),
					resource.TestCheckResourceAttr(resourceName1, "key_permissions.0", "create"),
					resource.TestCheckResourceAttr(resourceName1, "key_permissions.1", "get"),
					resource.TestCheckResourceAttr(resourceName1, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName1, "secret_permissions.1", "delete"),
					resource.TestCheckResourceAttr(resourceName1, "certificate_permissions.0", "create"),
					resource.TestCheckResourceAttr(resourceName1, "certificate_permissions.1", "delete"),

					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "key_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName2, "key_permissions.1", "encrypt"),
					resource.TestCheckResourceAttr(resourceName2, "secret_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName2, "secret_permissions.1", "delete"),
					resource.TestCheckResourceAttr(resourceName2, "certificate_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName2, "certificate_permissions.1", "delete"),
				),
			},
			{
				ResourceName:      resourceName1,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName2,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_update(t *testing.T) {
	rs := acctest.RandString(6)
	resourceName := "azurerm_key_vault_access_policy.test"
	preConfig := testAccAzureRMKeyVaultAccessPolicy_basic(rs, testLocation())
	postConfig := testAccAzureRMKeyVaultAccessPolicy_update(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "set"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.1", "encrypt"),
				),
			},
		},
	})
}

func testCheckAzureRMKeyVaultAccessPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		vaultName := rs.Primary.Attributes["vault_name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]
		objectId := rs.Primary.Attributes["object_id"]
		applicationId := rs.Primary.Attributes["application_id"]

		client := testAccProvider.Meta().(*ArmClient).keyVaultClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resGroup, vaultName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault %q (resource group: %q) does not exist", vaultName, resGroup)
			}

			return fmt.Errorf("Bad: Get on keyVaultClient: %+v", err)
		}

		policy, err := findKeyVaultAccessPolicy(resp.Properties.AccessPolicies, objectId, applicationId)
		if err != nil {
			return fmt.Errorf("Error finding Key Vault Access Policy %q : %+v", vaultName, err)
		}
		if policy == nil {
			return fmt.Errorf("Bad: Key Vault Policy %q (resource group: %q, object_id: %s) does not exist", vaultName, resGroup, objectId)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultAccessPolicy_basic(rString string, location string) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test" {
  vault_name          = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
    "set",
  ]

  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_multiple(rString string, location string) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test_with_application_id" {
  vault_name          = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  key_permissions = [
    "create",
    "get",
  ]

  secret_permissions = [
    "get",
    "delete",
  ]

  certificate_permissions = [
    "create",
    "delete",
  ]

  application_id = "${data.azurerm_client_config.current.service_principal_application_id}"
  tenant_id      = "${data.azurerm_client_config.current.tenant_id}"
  object_id      = "${data.azurerm_client_config.current.service_principal_object_id}"
}

resource "azurerm_key_vault_access_policy" "test_no_application_id" {
  vault_name          = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  key_permissions = [
    "list",
    "encrypt",
  ]

  secret_permissions = [
    "list",
    "delete",
  ]

  certificate_permissions = [
    "list",
    "delete",
  ]

  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_update(rString string, location string) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test" {
  vault_name          = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  key_permissions = [
    "list",
    "encrypt",
  ]

  secret_permissions = []

  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
}

`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_template(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  tags {
    environment = "Production"
  }
}
`, rString, location, rString)
}
