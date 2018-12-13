package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMManagementGroup_basic(t *testing.T) {
	resourceName := "azurerm_management_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMManagementGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(resourceName),
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

func TestAccAzureRMManagementGroup_nested(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMManagementGroup_nested(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists("azurerm_management_group.parent"),
					testCheckAzureRMManagementGroupExists("azurerm_management_group.child"),
				),
			},
			{
				ResourceName:      "azurerm_management_group.child",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMManagementGroup_multiLevel(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMManagementGroup_multiLevel(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists("azurerm_management_group.grandparent"),
					testCheckAzureRMManagementGroupExists("azurerm_management_group.parent"),
					testCheckAzureRMManagementGroupExists("azurerm_management_group.child"),
				),
			},
			{
				ResourceName:      "azurerm_management_group.child",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMManagementGroup_multiLevelUpdated(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMManagementGroup_nested(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists("azurerm_management_group.parent"),
					testCheckAzureRMManagementGroupExists("azurerm_management_group.child"),
				),
			},
			{
				Config: testAzureRMManagementGroup_multiLevel(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists("azurerm_management_group.grandparent"),
					testCheckAzureRMManagementGroupExists("azurerm_management_group.parent"),
					testCheckAzureRMManagementGroupExists("azurerm_management_group.child"),
				),
			},
		},
	})
}

func TestAccAzureRMManagementGroup_withName(t *testing.T) {
	resourceName := "azurerm_management_group.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMManagementGroup_withName(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMManagementGroup_updateName(t *testing.T) {
	resourceName := "azurerm_management_group.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMManagementGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(resourceName),
				),
			},
			{
				Config: testAzureRMManagementGroup_withName(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("acctestmg-%d", ri)),
				),
			},
		},
	})
}

func TestAccAzureRMManagementGroup_withSubscriptions(t *testing.T) {
	resourceName := "azurerm_management_group.test"
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMManagementGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subscription_ids.#", "0"),
				),
			},
			{
				Config: testAzureRMManagementGroup_withSubscriptions(subscriptionID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subscription_ids.#", "1"),
				),
			},
			{
				Config: testAzureRMManagementGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subscription_ids.#", "0"),
				),
			},
		},
	})
}

func testCheckAzureRMManagementGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		name := rs.Primary.Attributes["group_id"]

		client := testAccProvider.Meta().(*ArmClient).managementGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		recurse := false
		resp, err := client.Get(ctx, name, "", &recurse, "", "no-cache")
		if err != nil {
			return fmt.Errorf("Bad: Get on managementGroupsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Management Group does not exist: %s", name)
		}

		return nil
	}
}

func testCheckAzureRMManagementGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).managementGroupsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_management_group" {
			continue
		}

		name := rs.Primary.Attributes["group_id"]
		recurse := false
		resp, err := client.Get(ctx, name, "", &recurse, "", "no-cache")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Management Group still exists: %s", *resp.Name)
		}
	}

	return nil
}

func testAzureRMManagementGroup_basic() string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
}
`)
}

func testAzureRMManagementGroup_nested() string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "parent" {
}

resource "azurerm_management_group" "child" {
  parent_management_group_id = "${azurerm_management_group.parent.id}"
}
`)
}

func testAzureRMManagementGroup_multiLevel() string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "grandparent" {
}

resource "azurerm_management_group" "parent" {
  parent_management_group_id = "${azurerm_management_group.grandparent.id}"
}

resource "azurerm_management_group" "child" {
  parent_management_group_id = "${azurerm_management_group.parent.id}"
}
`)
}

func testAzureRMManagementGroup_withName(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%d"
}
`, rInt)
}

// TODO: switch this out for dynamically creating a subscription once that's supported in the future
func testAzureRMManagementGroup_withSubscriptions(subscriptionID string) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  subscription_ids = [
    "%s",
  ]
}
`, subscriptionID)
}
