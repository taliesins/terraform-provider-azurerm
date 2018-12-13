package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMApplicationInsights_basicWeb(t *testing.T) {
	resourceName := "azurerm_application_insights.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "web")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "application_type", "web"),
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

func TestAccAzureRMApplicationInsights_basicJava(t *testing.T) {
	resourceName := "azurerm_application_insights.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "java")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "application_type", "java"),
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

func TestAccAzureRMApplicationInsights_basicMobileCenter(t *testing.T) {
	resourceName := "azurerm_application_insights.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "MobileCenter")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "application_type", "MobileCenter"),
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

func TestAccAzureRMApplicationInsights_basicOther(t *testing.T) {
	resourceName := "azurerm_application_insights.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "other")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "application_type", "other"),
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

func TestAccAzureRMApplicationInsights_basicPhone(t *testing.T) {
	resourceName := "azurerm_application_insights.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "phone")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "application_type", "phone"),
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

func TestAccAzureRMApplicationInsights_basicStore(t *testing.T) {
	resourceName := "azurerm_application_insights.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "store")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "application_type", "store"),
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

func TestAccAzureRMApplicationInsights_basiciOS(t *testing.T) {
	resourceName := "azurerm_application_insights.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "ios")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "application_type", "ios"),
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

func testCheckAzureRMApplicationInsightsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).appInsightsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_insights" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Application Insights still exists:\n%#v", resp.ApplicationInsightsComponentProperties)
		}
	}

	return nil
}

func testCheckAzureRMApplicationInsightsExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Insights: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).appInsightsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on appInsightsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Application Insights '%q' (resource group: '%q') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMApplicationInsights_basic(rInt int, location string, applicationType string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "%s"
}
`, rInt, location, rInt, applicationType)
}
