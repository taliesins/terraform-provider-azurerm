---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_interface_backend_address_pool_association"
sidebar_current: "docs-azurerm-resource-network-interface-backend-address-pool-association"
description: |-
  Manages the association between a Network Interface and a Load Balancer's Backend Address Pool.

---

# azurerm_network_interface_backend_address_pool_association

Manages the association between a Network Interface and a Load Balancer's Backend Address Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "test" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "example-pip"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
  name                = "example-lb"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "primary"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "acctestpool"
}

resource "azurerm_network_interface" "test" {
  name                = "example-nic"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_network_interface_backend_address_pool_association" "test" {
  network_interface_id    = "${azurerm_network_interface.test.id}"
  ip_configuration_name   = "testconfiguration1"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_id` - (Required) The ID of the Network Interface. Changing this forces a new resource to be created.

* `ip_configuration_name` - (Required) The Name of the IP Configuration within the Network Interface which should be connected to the Backend Address Pool. Changing this forces a new resource to be created.

* `backend_address_pool_id` - (Required) The ID of the Load Balancer Backend Address Pool which this Network Interface which should be connected to. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The (Terraform specific) ID of the Association between the Network Interface and the Load Balancers Backend Address Pool.

## Import

Associations between Network Interfaces and Load Balancer Backend Address Pools can be imported using the `resource id`, e.g.


```shell
terraform import azurerm_network_interface_backend_address_pool_association.association1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.network/networkInterfaces/nic1/ipConfigurations/example|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/backendAddressPools/pool1
```

-> **NOTE:** This ID is specific to Terraform - and is of the format `{networkInterfaceId}/ipConfigurations/{ipConfigurationName}|{backendAddressPoolId}`.
