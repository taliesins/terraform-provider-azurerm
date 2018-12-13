package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAzureRMDatabrickWorkspaceName(t *testing.T) {
	cases := []struct {
		Value       string
		ShouldError bool
	}{
		{
			Value:       "hello",
			ShouldError: false,
		},
		{
			Value:       "hello123there",
			ShouldError: false,
		},
		{
			Value:       "hello-1-2-3-there",
			ShouldError: false,
		},
		{
			Value:       "hello-1-2-3-",
			ShouldError: true,
		},
		{
			Value:       "-hello-1-2-3",
			ShouldError: true,
		},
		{
			Value:       "hello!there",
			ShouldError: true,
		},
		{
			Value:       "hello--there",
			ShouldError: true,
		},
		{
			Value:       "!hellothere",
			ShouldError: true,
		},
		{
			Value:       "hellothere!",
			ShouldError: true,
		},
	}

	for _, tc := range cases {
		_, errors := validateDatabricksWorkspaceName(tc.Value, "test")

		hasErrors := len(errors) > 0
		if hasErrors && !tc.ShouldError {
			t.Fatalf("Expected no errors but got %d for %q", len(errors), tc.Value)
		}

		if !hasErrors && tc.ShouldError {
			t.Fatalf("Expected no errors but got %d for %q", len(errors), tc.Value)
		}
	}
}

func TestAccAzureRMDatabricksWorkspace_basic(t *testing.T) {
	resourceName := "azurerm_databricks_workspace.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDatabricksWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabricksWorkspace_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_resource_group_id"),
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

func TestAccAzureRMDatabricksWorkspace_complete(t *testing.T) {
	resourceName := "azurerm_databricks_workspace.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDatabricksWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabricksWorkspace_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_resource_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "managed_resource_group_name"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.pricing", "Standard"),
				),
			},
			{
				Config: testAccAzureRMDatabricksWorkspace_completeUpdate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_resource_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "managed_resource_group_name"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.pricing", "Standard"),
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

func testCheckAzureRMDatabricksWorkspaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Bad: Not found: %s", name)
		}

		workspaceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: No resource group found in state for Databricks Workspace: %s", workspaceName)
		}

		conn := testAccProvider.Meta().(*ArmClient).databricksWorkspacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, workspaceName)
		if err != nil {
			return fmt.Errorf("Bad: Getting Workspace: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Databricks Workspace %s (resource group: %s) does not exist", workspaceName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDatabricksWorkspaceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).databricksWorkspacesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_databricks_workspace" {
			continue
		}

		workspaceName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		resp, err := conn.Get(ctx, resourceGroup, workspaceName)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Bad: Databricks Workspace still exists:\n%#v", resp.ID)
		}
	}

	return nil
}

func testAccAzureRMDatabricksWorkspace_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
	name                = "acctestdbw-%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
	location            = "${azurerm_resource_group.test.location}"
	sku                 = "standard"
}
`, rInt, location, rInt)
}

func testAccAzureRMDatabricksWorkspace_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
	name                					= "acctestdbw-%d"
	resource_group_name 					= "${azurerm_resource_group.test.name}"
	location            					= "${azurerm_resource_group.test.location}"
	sku                 					= "standard"
	managed_resource_group_name 	= "acctestRG-%d-managed"

	tags {
		environment = "Production"
		pricing     = "Standard"
	}
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDatabricksWorkspace_completeUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
	name                					= "acctestdbw-%d"
	resource_group_name		 				= "${azurerm_resource_group.test.name}"
	location            					= "${azurerm_resource_group.test.location}"
	sku                 					= "standard"
	managed_resource_group_name 	= "acctestRG-%d-managed"

	tags {
		pricing = "Standard"
	}
}
`, rInt, location, rInt, rInt)
}
