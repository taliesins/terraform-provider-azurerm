package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultKey_basicEC(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_basicEC(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_size"},
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_basicRSA(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_basicRSA(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_size"},
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_basicRSAHSM(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_basicRSAHSM(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_size"},
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_complete(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_complete(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "world"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_size"},
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_update(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_basicRSA(rs, testLocation())
	updatedConfig := testAccAzureRMKeyVaultKey_basicUpdated(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_opts.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "key_opts.0", "decrypt"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_opts.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "key_opts.0", "encrypt"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_disappears(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_basicEC(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists(resourceName),
					testCheckAzureRMKeyVaultKeyDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_disappearsWhenParentKeyVaultDeleted(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_basicEC(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultKeyExists("azurerm_key_vault_key.test"),
					testCheckAzureRMKeyVaultDisappears("azurerm_key_vault.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMKeyVaultKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).keyVaultManagementClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault_key" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]

		// get the latest version
		resp, err := client.GetKey(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Key Vault Key still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMKeyVaultKeyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]

		client := testAccProvider.Meta().(*ArmClient).keyVaultManagementClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.GetKey(ctx, vaultBaseUrl, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault Key %q (resource group: %q) does not exist", name, vaultBaseUrl)
			}

			return fmt.Errorf("Bad: Get on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMKeyVaultKeyDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		vaultBaseUrl := rs.Primary.Attributes["vault_uri"]

		client := testAccProvider.Meta().(*ArmClient).keyVaultManagementClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.DeleteKey(ctx, vaultBaseUrl, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Delete on keyVaultManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultKey_basicEC(rString string, location string) string {
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

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
      "delete",
      "get",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "EC"
  key_size  = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}
`, rString, location, rString, rString)
}

func testAccAzureRMKeyVaultKey_basicRSA(rString string, location string) string {
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

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "RSA"
  key_size  = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}
`, rString, location, rString, rString)
}

func testAccAzureRMKeyVaultKey_basicRSAHSM(rString string, location string) string {
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

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
      "delete",
      "get",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "RSA-HSM"
  key_size  = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}
`, rString, location, rString, rString)
}

func testAccAzureRMKeyVaultKey_complete(rString string, location string) string {
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

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
      "delete",
      "get",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "RSA"
  key_size  = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  tags {
    "hello" = "world"
  }
}
`, rString, location, rString, rString)
}

func testAccAzureRMKeyVaultKey_basicUpdated(rString string, location string) string {
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

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "test" {
  name      = "key-%s"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "RSA"
  key_size  = 2048

  key_opts = [
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}
`, rString, location, rString, rString)
}
