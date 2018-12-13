---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet_route_table_association"
sidebar_current: "docs-azurerm-resource-network-subnet-route-table-association"
description: |-
  Associates a [Route Table](route_table.html) with a [Subnet](subnet.html) within a [Virtual Network](virtual_network.html).

---

# azurerm_subnet_route_table_association

Associates a [Route Table](route_table.html) with a [Subnet](subnet.html) within a [Virtual Network](virtual_network.html).

-> **NOTE:** Subnet `<->` Route Table associations currently need to be configured on both this resource and using the `route_table_id` field on the `azurerm_subnet` resource. The next major version of the AzureRM Provider (2.0) will remove the `route_table_id` field from the `azurerm_subnet` resource such that this resource is used to link resources in future.

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
  name                 = "frontend"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  route_table_id       = "${azurerm_route_table.test.id}"
}

resource "azurerm_route_table" "test" {
  name                = "example-routetable"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name                   = "example"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = "${azurerm_subnet.test.id}"
  route_table_id = "${azurerm_route_table.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required) The ID of the Route Table which should be associated with the Subnet. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Subnet.

## Import

Subnet Route Table Associations can be imported using the `resource id` of the Subnet, e.g.

```shell
terraform import azurerm_subnet_route_table_association.association1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```
