package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMLoadBalancerBackEndAddressPool_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	addressPoolName := fmt.Sprintf("%d-address-pool", ri)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	backendAddressPoolId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/backendAddressPools/%s",
		subscriptionID, ri, ri, addressPoolName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(ri, addressPoolName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_backend_address_pool.test", "id", backendAddressPoolId),
				),
			},
			{
				ResourceName:      "azurerm_lb.test",
				ImportState:       true,
				ImportStateVerify: true,
				// location is deprecated and was never actually used
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_removal(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	addressPoolName := fmt.Sprintf("%d-address-pool", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_removal(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolNotExists(addressPoolName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_reapply(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	addressPoolName := fmt.Sprintf("%d-address-pool", ri)

	deleteAddressPoolState := func(s *terraform.State) error {
		return s.Remove("azurerm_lb_backend_address_pool.test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(ri, addressPoolName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName, &lb),
					deleteAddressPoolState,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(ri, addressPoolName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_disappears(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	addressPoolName := fmt.Sprintf("%d-address-pool", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(ri, addressPoolName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName, &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolDisappears(addressPoolName, &lb),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerBackEndAddressPoolByName(lb, addressPoolName)
		if !exists {
			return fmt.Errorf("A BackEnd Address Pool with name %q cannot be found.", addressPoolName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerBackEndAddressPoolNotExists(addressPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerBackEndAddressPoolByName(lb, addressPoolName)
		if exists {
			return fmt.Errorf("A BackEnd Address Pool with name %q has been found.", addressPoolName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerBackEndAddressPoolDisappears(addressPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).loadBalancerClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, i, exists := findLoadBalancerBackEndAddressPoolByName(lb, addressPoolName)
		if !exists {
			return fmt.Errorf("A BackEnd Address Pool with name %q cannot be found.", addressPoolName)
		}

		currentPools := *lb.LoadBalancerPropertiesFormat.BackendAddressPools
		pools := append(currentPools[:i], currentPools[i+1:]...)
		lb.LoadBalancerPropertiesFormat.BackendAddressPools = &pools

		id, err := parseAzureResourceID(*lb.ID)
		if err != nil {
			return err
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, *lb.Name, *lb)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating Load Balancer %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error Creating/Updating Load Balancer %+v", err)
		}

		_, err = client.Get(ctx, id.ResourceGroup, *lb.Name, "")
		return err
	}
}

func testAccAzureRMLoadBalancerBackEndAddressPool_basic(rInt int, addressPoolName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                         = "test-ip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
}
`, rInt, location, rInt, rInt, rInt, addressPoolName)
}

func testAccAzureRMLoadBalancerBackEndAddressPool_removal(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                         = "test-ip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}
