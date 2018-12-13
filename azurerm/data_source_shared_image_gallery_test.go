package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMSharedImageGallery_basic(t *testing.T) {
	dataSourceName := "data.azurerm_shared_image_gallery.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageGalleryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedImageGallery_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSharedImageGallery_complete(t *testing.T) {
	dataSourceName := "data.azurerm_shared_image_gallery.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageGalleryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedImageGallery_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "description", "Shared images and things."),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.Hello", "There"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.World", "Example"),
				),
			},
		},
	})
}

func testAccDataSourceSharedImageGallery_basic(rInt int, location string) string {
	template := testAccAzureRMSharedImageGallery_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_gallery" "test" {
  name                = "${azurerm_shared_image_gallery.test.name}"
  resource_group_name = "${azurerm_shared_image_gallery.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceSharedImageGallery_complete(rInt int, location string) string {
	template := testAccAzureRMSharedImageGallery_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_gallery" "test" {
  name                = "${azurerm_shared_image_gallery.test.name}"
  resource_group_name = "${azurerm_shared_image_gallery.test.resource_group_name}"
}
`, template)
}
