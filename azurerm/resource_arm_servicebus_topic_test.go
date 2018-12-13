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

func TestAccAzureRMServiceBusTopic_basic(t *testing.T) {
	resourceName := "azurerm_servicebus_topic.test"
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopic_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
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

func TestAccAzureRMServiceBusTopic_basicDisabled(t *testing.T) {
	resourceName := "azurerm_servicebus_topic.test"
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopic_basicDisabled(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
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

func TestAccAzureRMServiceBusTopic_basicDisableEnable(t *testing.T) {
	resourceName := "azurerm_servicebus_topic.test"
	ri := acctest.RandInt()
	location := testLocation()
	enabledConfig := testAccAzureRMServiceBusTopic_basic(ri, location)
	disabledConfig := testAccAzureRMServiceBusTopic_basicDisabled(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: enabledConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
				),
			},
			{
				Config: disabledConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
				),
			},
			{
				Config: enabledConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusTopic_update(t *testing.T) {
	resourceName := "azurerm_servicebus_topic.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusTopic_basic(ri, location)
	postConfig := testAccAzureRMServiceBusTopic_update(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_batched_operations", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_express", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusTopic_enablePartitioningStandard(t *testing.T) {
	resourceName := "azurerm_servicebus_topic.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusTopic_basic(ri, location)
	postConfig := testAccAzureRMServiceBusTopic_enablePartitioningStandard(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_partitioning", "true"),
					// Ensure size is read back in it's original value and not the x16 value returned by Azure
					resource.TestCheckResourceAttr(resourceName, "max_size_in_megabytes", "5120"),
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

func TestAccAzureRMServiceBusTopic_enablePartitioningPremium(t *testing.T) {
	resourceName := "azurerm_servicebus_topic.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusTopic_basicPremium(ri, location)
	postConfig := testAccAzureRMServiceBusTopic_enablePartitioningPremium(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_partitioning", "false"),
					resource.TestCheckResourceAttr(resourceName, "max_size_in_megabytes", "81920"),
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

func TestAccAzureRMServiceBusTopic_enableDuplicateDetection(t *testing.T) {
	resourceName := "azurerm_servicebus_topic.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusTopic_basic(ri, location)
	postConfig := testAccAzureRMServiceBusTopic_enableDuplicateDetection(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "requires_duplicate_detection", "true"),
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

func TestAccAzureRMServiceBusTopic_isoTimeSpanAttributes(t *testing.T) {
	resourceName := "azurerm_servicebus_topic.test"
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopic_isoTimeSpanAttributes(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_delete_on_idle", "PT10M"),
					resource.TestCheckResourceAttr(resourceName, "default_message_ttl", "PT30M"),
					resource.TestCheckResourceAttr(resourceName, "requires_duplicate_detection", "true"),
					resource.TestCheckResourceAttr(resourceName, "duplicate_detection_history_time_window", "PT15M"),
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

func testCheckAzureRMServiceBusTopicDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).serviceBusTopicsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_topic" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("ServiceBus Topic still exists:\n%+v", resp.SBTopicProperties)
		}
	}

	return nil
}

func testCheckAzureRMServiceBusTopicExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		topicName := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for topic: %s", topicName)
		}

		client := testAccProvider.Meta().(*ArmClient).serviceBusTopicsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceBusTopicsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Topic %q (resource group: %q) does not exist", namespaceName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceBusTopic_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusTopic_basicDisabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  status              = "disabled"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusTopic_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                      = "acctestservicebustopic-%d"
  namespace_name            = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  enable_batched_operations = true
  enable_express            = true
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusTopic_basicPremium(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Premium"
  capacity            = 1
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  enable_partitioning = false
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusTopic_enablePartitioningStandard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                  = "acctestservicebustopic-%d"
  namespace_name        = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  enable_partitioning   = true
  max_size_in_megabytes = 5120
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusTopic_enablePartitioningPremium(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "premium"
  capacity            = 1
}

resource "azurerm_servicebus_topic" "test" {
  name                  = "acctestservicebustopic-%d"
  namespace_name        = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  enable_partitioning   = false
  max_size_in_megabytes = 81920
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusTopic_enableDuplicateDetection(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                         = "acctestservicebustopic-%d"
  namespace_name               = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  requires_duplicate_detection = true
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusTopic_isoTimeSpanAttributes(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                                    = "acctestservicebustopic-%d"
  namespace_name                          = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name                     = "${azurerm_resource_group.test.name}"
  auto_delete_on_idle                     = "PT10M"
  default_message_ttl                     = "PT30M"
  requires_duplicate_detection            = true
  duplicate_detection_history_time_window = "PT15M"
}
`, rInt, location, rInt, rInt)
}
