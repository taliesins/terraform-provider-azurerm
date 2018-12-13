package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDevTestVirtualMachine_basic(t *testing.T) {
	resourceName := "azurerm_dev_test_windows_virtual_machine.test"
	rInt := acctest.RandIntRange(11111, 99999)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevTestWindowsVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestWindowsVirtualMachine_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestWindowsVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gallery_image_reference.0.publisher", "MicrosoftWindowsServer"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"lab_subnet_name",
					"lab_virtual_network_id",
					"password",
				},
			},
		},
	})
}

func TestAccAzureRMDevTestWindowsVirtualMachine_inboundNatRules(t *testing.T) {
	resourceName := "azurerm_dev_test_windows_virtual_machine.test"
	rInt := acctest.RandIntRange(11111, 99999)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevTestWindowsVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestWindowsVirtualMachine_inboundNatRules(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestWindowsVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "disallow_public_ip_address", "true"),
					resource.TestCheckResourceAttr(resourceName, "gallery_image_reference.0.publisher", "MicrosoftWindowsServer"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Acceptance", "Test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"inbound_nat_rule",
					"lab_subnet_name",
					"lab_virtual_network_id",
					"password",
				},
			},
		},
	})
}

func TestAccAzureRMDevTestWindowsVirtualMachine_updateStorage(t *testing.T) {
	resourceName := "azurerm_dev_test_windows_virtual_machine.test"
	rInt := acctest.RandIntRange(11111, 99999)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevTestWindowsVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestWindowsVirtualMachine_storage(rInt, location, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestWindowsVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gallery_image_reference.0.publisher", "MicrosoftWindowsServer"),
					resource.TestCheckResourceAttr(resourceName, "storage_type", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMDevTestWindowsVirtualMachine_storage(rInt, location, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestWindowsVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gallery_image_reference.0.publisher", "MicrosoftWindowsServer"),
					resource.TestCheckResourceAttr(resourceName, "storage_type", "Premium"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testCheckAzureRMDevTestWindowsVirtualMachineExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		virtualMachineName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).devTestVirtualMachinesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualMachineName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get devTestVirtualMachinesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevTest Windows Virtual Machine %q (Lab %q / Resource Group: %q) does not exist", virtualMachineName, labName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDevTestWindowsVirtualMachineDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).devTestVirtualMachinesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_windows_virtual_machine" {
			continue
		}

		virtualMachineName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualMachineName, "")

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DevTest Windows Virtual Machine still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDevTestWindowsVirtualMachine_basic(rInt int, location string) string {
	template := testAccAzureRMDevTestWindowsVirtualMachine_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "test" {
  name                   = "acctestvm%d"
  lab_name               = "${azurerm_dev_test_lab.test.name}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  location               = "${azurerm_resource_group.test.location}"
  size                   = "Standard_F2"
  username               = "acct5stU5er"
  password               = "Pa$$w0rd1234!"
  lab_virtual_network_id = "${azurerm_dev_test_virtual_network.test.id}"
  lab_subnet_name        = "${azurerm_dev_test_virtual_network.test.subnet.0.name}"
  storage_type           = "Standard"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }
}
`, template, rInt)
}

func testAccAzureRMDevTestWindowsVirtualMachine_inboundNatRules(rInt int, location string) string {
	template := testAccAzureRMDevTestWindowsVirtualMachine_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "test" {
  name                       = "acctestvm%d"
  lab_name                   = "${azurerm_dev_test_lab.test.name}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
  location                   = "${azurerm_resource_group.test.location}"
  size                       = "Standard_F2"
  username                   = "acct5stU5er"
  password                   = "Pa$$w0rd1234!"
  disallow_public_ip_address = true
  lab_virtual_network_id     = "${azurerm_dev_test_virtual_network.test.id}"
  lab_subnet_name            = "${azurerm_dev_test_virtual_network.test.subnet.0.name}"
  storage_type               = "Standard"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }

  inbound_nat_rule {
    protocol     = "Tcp"
    backend_port = 22
  }

  inbound_nat_rule {
    protocol     = "Tcp"
    backend_port = 3389
  }

  tags {
    "Acceptance" = "Test"
  }
}
`, template, rInt)
}

func testAccAzureRMDevTestWindowsVirtualMachine_storage(rInt int, location, storageType string) string {
	template := testAccAzureRMDevTestWindowsVirtualMachine_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "test" {
  name                   = "acctestvm%d"
  lab_name               = "${azurerm_dev_test_lab.test.name}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  location               = "${azurerm_resource_group.test.location}"
  size                   = "Standard_B1ms"
  username               = "acct5stU5er"
  password               = "Pa$$w0rd1234!"
  lab_virtual_network_id = "${azurerm_dev_test_virtual_network.test.id}"
  lab_subnet_name        = "${azurerm_dev_test_virtual_network.test.subnet.0.name}"
  storage_type           = "%s"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }
}
`, template, rInt, storageType)
}

func testAccAzureRMDevTestWindowsVirtualMachine_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dev_test_virtual_network" "test" {
  name                = "acctestdtvn%d"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  subnet {
    use_public_ip_address           = "Allow"
    use_in_virtual_machine_creation = "Allow"
  }
}
`, rInt, location, rInt, rInt)
}
