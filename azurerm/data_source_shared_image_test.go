package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMSharedImage_basic(t *testing.T) {
	dataSourceName := "data.azurerm_shared_image.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedImage_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSharedImage_complete(t *testing.T) {
	dataSourceName := "data.azurerm_shared_image.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedImage_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccDataSourceSharedImage_basic(rInt int, location string) string {
	template := testAccAzureRMSharedImage_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = "${azurerm_shared_image.test.name}"
  gallery_name        = "${azurerm_shared_image.test.gallery_name}"
  resource_group_name = "${azurerm_shared_image.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceSharedImage_complete(rInt int, location string) string {
	template := testAccAzureRMSharedImage_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = "${azurerm_shared_image.test.name}"
  gallery_name        = "${azurerm_shared_image.test.gallery_name}"
  resource_group_name = "${azurerm_shared_image.test.resource_group_name}"
}
`, template)
}
