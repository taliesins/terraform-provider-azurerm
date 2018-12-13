package azurerm

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceArmStorageAccountSas_basic(t *testing.T) {
	dataSourceName := "data.azurerm_storage_account_sas.test"
	rInt := acctest.RandInt()
	rString := acctest.RandString(4)
	location := testLocation()
	utcNow := time.Now().UTC()
	startDate := utcNow.Format(time.RFC3339)
	endDate := utcNow.Add(time.Hour * 24).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccountSas_basic(rInt, rString, location, startDate, endDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "https_only", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "start", startDate),
					resource.TestCheckResourceAttr(dataSourceName, "expiry", endDate),
					resource.TestCheckResourceAttrSet(dataSourceName, "sas"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageAccountSas_basic(rInt int, rString string, location string, startDate string, endDate string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestsa-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsads%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "production"
  }
}

data "azurerm_storage_account_sas" "test" {
  connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  https_only        = true

  resource_types {
    service   = true
    container = false
    object    = false
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "%s"
  expiry = "%s"

  permissions {
    read    = true
    write   = true
    delete  = false
    list    = false
    add     = true
    create  = true
    update  = false
    process = false
  }
}
`, rInt, location, rString, startDate, endDate)
}

func TestAccDataSourceArmStorageAccountSas_resourceTypesString(t *testing.T) {
	testCases := []struct {
		input    map[string]interface{}
		expected string
	}{
		{map[string]interface{}{"service": true}, "s"},
		{map[string]interface{}{"container": true}, "c"},
		{map[string]interface{}{"object": true}, "o"},
		{map[string]interface{}{"service": true, "container": true, "object": true}, "sco"},
	}

	for _, test := range testCases {
		result := buildResourceTypesString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}

func TestAccDataSourceArmStorageAccountSas_servicesString(t *testing.T) {
	testCases := []struct {
		input    map[string]interface{}
		expected string
	}{
		{map[string]interface{}{"blob": true}, "b"},
		{map[string]interface{}{"queue": true}, "q"},
		{map[string]interface{}{"table": true}, "t"},
		{map[string]interface{}{"file": true}, "f"},
		{map[string]interface{}{"blob": true, "queue": true, "table": true, "file": true}, "bqtf"},
	}

	for _, test := range testCases {
		result := buildServicesString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}

}

func TestAccDataSourceArmStorageAccountSas_permissionsString(t *testing.T) {
	testCases := []struct {
		input    map[string]interface{}
		expected string
	}{
		{map[string]interface{}{"read": true}, "r"},
		{map[string]interface{}{"write": true}, "w"},
		{map[string]interface{}{"delete": true}, "d"},
		{map[string]interface{}{"list": true}, "l"},
		{map[string]interface{}{"add": true}, "a"},
		{map[string]interface{}{"create": true}, "c"},
		{map[string]interface{}{"update": true}, "u"},
		{map[string]interface{}{"process": true}, "p"},
		{map[string]interface{}{"read": true, "write": true, "add": true, "create": true}, "rwac"},
	}

	for _, test := range testCases {
		result := buildPermissionsString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}
