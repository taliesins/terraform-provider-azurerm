package azurerm

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationSchedule_oneTime_basic(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTime_basic(ri, testLocation()),
				Check:  checkAccAzureRMAutomationSchedule_oneTime_basic(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_oneTime_complete(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	//the API returns the time in the timezone we pass in
	//it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("CET")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTime_complete(ri, testLocation(), startTime),
				Check:  checkAccAzureRMAutomationSchedule_oneTime_complete(resourceName, startTime),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_oneTime_update(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	//the API returns the time in the timezone we pass in
	//it also seems to strip seconds, hijack the RFC3339 format to have 0s there
	loc, _ := time.LoadLocation("CET")
	startTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_oneTime_basic(ri, testLocation()),
				Check:  checkAccAzureRMAutomationSchedule_oneTime_basic(resourceName),
			},
			{
				Config: testAccAzureRMAutomationSchedule_oneTime_complete(ri, testLocation(), startTime),
				Check:  checkAccAzureRMAutomationSchedule_oneTime_complete(resourceName, startTime),
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_hourly(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_basic(ri, testLocation(), "Hour", 7),
				Check:  checkAccAzureRMAutomationSchedule_recurring_basic(resourceName, "Hour", 7),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_daily(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_basic(ri, testLocation(), "Day", 7),
				Check:  checkAccAzureRMAutomationSchedule_recurring_basic(resourceName, "Day", 7),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_weekly(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_basic(ri, testLocation(), "Week", 7),
				Check:  checkAccAzureRMAutomationSchedule_recurring_basic(resourceName, "Week", 7),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_monthly(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_basic(ri, testLocation(), "Month", 7),
				Check:  checkAccAzureRMAutomationSchedule_recurring_basic(resourceName, "Month", 7),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_weekly_advanced(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_advanced_week(ri, testLocation(), "Monday"),
				Check:  checkAccAzureRMAutomationSchedule_recurring_advanced_week(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_monthly_advanced_by_day(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_advanced_month(ri, testLocation(), 2),
				Check:  checkAccAzureRMAutomationSchedule_recurring_advanced_month(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationSchedule_monthly_advanced_by_week_day(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationSchedule_recurring_advanced_month_week_day(ri, testLocation(), "Monday", 2),
				Check:  checkAccAzureRMAutomationSchedule_recurring_advanced_month_week_day(resourceName, "Monday", 2),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMAutomationScheduleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationScheduleClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_schedule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Schedule: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Schedule still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationScheduleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*ArmClient).automationScheduleClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Schedule: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Schedule '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationScheduleClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAutomationSchedule_prerequisites(rInt int, location string) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAutomationSchedule_oneTime_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "OneTime"
}
`, testAccAzureRMAutomationSchedule_prerequisites(rInt, location), rInt)
}

func checkAccAzureRMAutomationSchedule_oneTime_basic(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "OneTime"),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
	)
}

func testAccAzureRMAutomationSchedule_oneTime_complete(rInt int, location, startTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "OneTime"
  start_time              = "%s"
  timezone                = "Central Europe Standard Time" 
  description             = "This is an automation schedule"
}
`, testAccAzureRMAutomationSchedule_prerequisites(rInt, location), rInt, startTime)
}

func checkAccAzureRMAutomationSchedule_oneTime_complete(resourceName, startTime string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "OneTime"),
		resource.TestCheckResourceAttr(resourceName, "start_time", startTime),
		resource.TestCheckResourceAttr(resourceName, "timezone", "Central Europe Standard Time"),
		resource.TestCheckResourceAttr(resourceName, "description", "This is an automation schedule"),
	)
}

func testAccAzureRMAutomationSchedule_recurring_basic(rInt int, location, frequency string, interval int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "%s"
  interval                = "%d"
}
`, testAccAzureRMAutomationSchedule_prerequisites(rInt, location), rInt, frequency, interval)
}

func checkAccAzureRMAutomationSchedule_recurring_basic(resourceName string, frequency string, interval int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", frequency),
		resource.TestCheckResourceAttr(resourceName, "interval", strconv.Itoa(interval)),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
	)
}

func testAccAzureRMAutomationSchedule_recurring_advanced_week(rInt int, location string, weekDay string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "Week"
  interval                = "1"
  week_days               = ["%s"]
}	
`, testAccAzureRMAutomationSchedule_prerequisites(rInt, location), rInt, weekDay)
}

func checkAccAzureRMAutomationSchedule_recurring_advanced_week(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists("azurerm_automation_schedule.test"),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "Week"),
		resource.TestCheckResourceAttr(resourceName, "interval", "1"),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
		resource.TestCheckResourceAttr(resourceName, "week_days.#", "1"),
	)
}

func testAccAzureRMAutomationSchedule_recurring_advanced_month(rInt int, location string, monthDay int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "Month"
  interval                = "1"
  month_days              = [%d]
}	
`, testAccAzureRMAutomationSchedule_prerequisites(rInt, location), rInt, monthDay)
}

func checkAccAzureRMAutomationSchedule_recurring_advanced_month(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists("azurerm_automation_schedule.test"),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "Month"),
		resource.TestCheckResourceAttr(resourceName, "interval", "1"),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
		resource.TestCheckResourceAttr(resourceName, "month_days.#", "1"),
	)
}

func testAccAzureRMAutomationSchedule_recurring_advanced_month_week_day(rInt int, location string, weekDay string, weekDayOccurrence int) string {

	return fmt.Sprintf(`
%s

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "Month"
  interval                = "1"

  monthly_occurrence {
	day        = "%s"
	occurrence = "%d"
  }
}	
`, testAccAzureRMAutomationSchedule_prerequisites(rInt, location), rInt, weekDay, weekDayOccurrence)
}

func checkAccAzureRMAutomationSchedule_recurring_advanced_month_week_day(resourceName string, monthWeekDay string, monthWeekOccurrence int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationScheduleExists("azurerm_automation_schedule.test"),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "start_time"),
		resource.TestCheckResourceAttr(resourceName, "frequency", "Month"),
		resource.TestCheckResourceAttr(resourceName, "interval", "1"),
		resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
		resource.TestCheckResourceAttr(resourceName, "monthly_occurrence.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "monthly_occurrence.0.day", monthWeekDay),
		resource.TestCheckResourceAttr(resourceName, "monthly_occurrence.0.occurrence", strconv.Itoa(monthWeekOccurrence)),
	)
}
