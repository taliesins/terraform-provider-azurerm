package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCognitiveAccount_basic(t *testing.T) {
	resourceName := "azurerm_cognitive_account.test"
	ri := acctest.RandInt()
	config := testAccAzureRMCognitiveAccount_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccAzureRMCognitiveAccount_complete(t *testing.T) {
	resourceName := "azurerm_cognitive_account.test"
	ri := acctest.RandInt()
	config := testAccAzureRMCognitiveAccount_complete(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Acceptance", "Test"),
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

func TestAccAzureRMCognitiveAccount_update(t *testing.T) {
	resourceName := "azurerm_cognitive_account.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCognitiveAccount_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMCognitiveAccount_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Acceptance", "Test"),
				),
			},
		},
	})
}

func testCheckAzureRMAppCognitiveAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).cognitiveAccountsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cognitive_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetProperties(ctx, resourceGroup, name)
		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("Cognitive Services Account still exists:\n%#v", resp)
			}

			return nil
		}

	}

	return nil
}

func testCheckAzureRMCognitiveAccountExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).cognitiveAccountsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.GetProperties(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Cognitive Services Account %q (Resource Group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on cognitiveAccountsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMCognitiveAccount_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Face"

  sku {
    name = "S0"
    tier = "Standard"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCognitiveAccount_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Face"

  sku {
    name = "S0"
    tier = "Standard"
  }

  tags {
    Acceptance = "Test"
  }
}
`, rInt, location, rInt)
}
