package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_europe(t *testing.T) {
	dataSourceName := "data.azurerm_traffic_manager_geographical_location.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template("Europe"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "GEO-EU"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "Europe"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_germany(t *testing.T) {
	dataSourceName := "data.azurerm_traffic_manager_geographical_location.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template("Germany"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "DE"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "Germany"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_unitedKingdom(t *testing.T) {
	dataSourceName := "data.azurerm_traffic_manager_geographical_location.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template("United Kingdom"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "GB"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "United Kingdom"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceTrafficManagerGeographicalLocation_world(t *testing.T) {
	dataSourceName := "data.azurerm_traffic_manager_geographical_location.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template("World"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "WORLD"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "World"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceTrafficManagerGeographicalLocation_template(name string) string {
	return fmt.Sprintf(`
data "azurerm_traffic_manager_geographical_location" "test" {
  name = "%s"
}
`, name)
}
