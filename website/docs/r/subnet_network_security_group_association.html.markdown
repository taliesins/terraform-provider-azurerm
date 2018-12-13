---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet_network_security_group_association"
sidebar_current: "docs-azurerm-resource-network-subnet-network-security-group-association"
description: |-
  Associates a [Network Security Group](network_security_group.html) with a [Subnet](subnet.html) within a [Virtual Network](virtual_network.html).

---

# azurerm_subnet_network_security_group_association

Associates a [Network Security Group](network_security_group.html) with a [Subnet](subnet.html) within a [Virtual Network](virtual_network.html).

-> **NOTE:** Subnet `<->` Network Security Group associations currently need to be configured on both this resource and using the `network_security_group_id` field on the `azurerm_subnet` resource. The next major version of the AzureRM Provider (2.0) will remove the `network_security_group_id` field from the `azurerm_subnet` resource such that this resource is used to link resources in future.

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
  name                      = "frontend"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test.name}"
  address_prefix            = "10.0.2.0/24"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}

resource "azurerm_network_security_group" "test" {
  name                = "example-nsg"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = "${azurerm_subnet.test.id}"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `network_security_group_id` - (Required) The ID of the Network Security Group which should be associated with the Subnet. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Subnet.

## Import

Subnet `<->` Network Security Group Associations can be imported using the `resource id` of the Subnet, e.g.

```shell
terraform import azurerm_subnet_network_security_group_association.association1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```
