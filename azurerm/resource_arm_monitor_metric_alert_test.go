package azurerm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMMonitorMetricAlert_basic(t *testing.T) {
	resourceName := "azurerm_monitor_metric_alert.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMMonitorMetricAlert_basic(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aggregation", "Average"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operator", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.threshold", "55.5"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "0"),
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

func TestAccAzureRMMonitorMetricAlert_complete(t *testing.T) {
	resourceName := "azurerm_monitor_metric_alert.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMMonitorMetricAlert_complete(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_mitigate", "true"),
					resource.TestCheckResourceAttr(resourceName, "severity", "4"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a complete metric alert resource."),
					resource.TestCheckResourceAttr(resourceName, "frequency", "PT30M"),
					resource.TestCheckResourceAttr(resourceName, "window_size", "PT12H"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_name", "Transactions"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aggregation", "Maximum"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operator", "Equals"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.threshold", "99"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.0.name", "GeoType"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.0.operator", "Include"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.0.values.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.0.values.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.1.name", "ApiName"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.1.operator", "Include"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.1.values.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.1.values.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.aggregation", "Total"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.operator", "GreaterThanOrEqual"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.threshold", "66.6"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.dimension.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
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

func TestAccAzureRMMonitorMetricAlert_basicAndCompleteUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_metric_alert.test"
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	basicConfig := testAccAzureRMMonitorMetricAlert_basic(ri, rs, location)
	completeConfig := testAccAzureRMMonitorMetricAlert_complete(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_mitigate", "false"),
					resource.TestCheckResourceAttr(resourceName, "severity", "3"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "frequency", "PT1M"),
					resource.TestCheckResourceAttr(resourceName, "window_size", "PT5M"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aggregation", "Average"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operator", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.threshold", "55.5"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "0"),
				),
			},
			{
				Config: completeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_mitigate", "true"),
					resource.TestCheckResourceAttr(resourceName, "severity", "4"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a complete metric alert resource."),
					resource.TestCheckResourceAttr(resourceName, "frequency", "PT30M"),
					resource.TestCheckResourceAttr(resourceName, "window_size", "PT12H"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_name", "Transactions"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aggregation", "Maximum"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operator", "Equals"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.threshold", "99"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.0.name", "GeoType"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.0.operator", "Include"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.0.values.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.0.values.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.1.name", "ApiName"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.1.operator", "Include"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.1.values.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.1.values.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.aggregation", "Total"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.operator", "GreaterThanOrEqual"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.threshold", "66.6"),
					resource.TestCheckResourceAttr(resourceName, "criteria.1.dimension.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
				),
			},
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_mitigate", "false"),
					resource.TestCheckResourceAttr(resourceName, "severity", "3"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "frequency", "PT1M"),
					resource.TestCheckResourceAttr(resourceName, "window_size", "PT5M"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aggregation", "Average"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operator", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.threshold", "55.5"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.dimension.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "0"),
				),
			},
		},
	})
}

func testAccAzureRMMonitorMetricAlert_basic(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  scopes              = ["${azurerm_storage_account.test.id}"]

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "UsedCapacity"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = 55.5
  }
}
`, rInt, location, rString, rInt)
}

func testAccAzureRMMonitorMetricAlert_complete(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa1%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_action_group" "test1" {
  name                = "acctestActionGroup1-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag1"
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag2"
}

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  scopes              = ["${azurerm_storage_account.test.id}"]
  enabled             = true
  auto_mitigate       = true
  severity            = 4
  description         = "This is a complete metric alert resource."
  frequency           = "PT30M"
  window_size         = "PT12H"

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "Transactions"
    aggregation      = "Maximum"
    operator         = "Equals"
    threshold        = 99

    dimension {
      "name"     = "GeoType"
      "operator" = "Include"
      "values"   = ["*"]
    }

    dimension {
      "name"     = "ApiName"
      "operator" = "Include"
      "values"   = ["*"]
    }
  }

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "UsedCapacity"
    aggregation      = "Total"
    operator         = "GreaterThanOrEqual"
    threshold        = 66.6
  }

  action {
    action_group_id = "${azurerm_monitor_action_group.test1.id}"
  }

  action {
    action_group_id = "${azurerm_monitor_action_group.test2.id}"
  }
}
`, rInt, location, rString, rInt, rInt, rInt)
}

func testCheckAzureRMMonitorMetricAlertDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).monitorMetricAlertsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_metric_alert" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Metric alert still exists:\n%#v", resp)
		}
	}
	return nil
}

func testCheckAzureRMMonitorMetricAlertExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Metric Alert Instance: %s", resourceName)
		}

		conn := testAccProvider.Meta().(*ArmClient).monitorMetricAlertsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, resourceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on monitorMetricAlertsClient: %+v", err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Metric Alert Instance %q (resource group: %q) does not exist", resourceName, resourceGroup)
		}

		return nil
	}
}
