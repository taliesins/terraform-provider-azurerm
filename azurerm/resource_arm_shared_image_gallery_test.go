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

func TestAccAzureRMSharedImageGallery_basic(t *testing.T) {
	resourceName := "azurerm_shared_image_gallery.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageGalleryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImageGallery_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageGalleryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func TestAccAzureRMSharedImageGallery_complete(t *testing.T) {
	resourceName := "azurerm_shared_image_gallery.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageGalleryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImageGallery_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageGalleryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "Shared images and things."),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "There"),
					resource.TestCheckResourceAttr(resourceName, "tags.World", "Example"),
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

func testCheckAzureRMSharedImageGalleryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).galleriesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_shared_image_gallery" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Shared Image Gallery still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMSharedImageGalleryExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		galleryName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Shared Image Gallery: %s", galleryName)
		}

		client := testAccProvider.Meta().(*ArmClient).galleriesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, galleryName)
		if err != nil {
			return fmt.Errorf("Bad: Get on galleriesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Shared Image Gallery %q (resource group: %q) does not exist", galleryName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMSharedImageGallery_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, rInt, location, rInt)
}

func testAccAzureRMSharedImageGallery_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  description         = "Shared images and things."

  tags {
    Hello = "There"
    World = "Example"
  }
}
`, rInt, location, rInt)
}
