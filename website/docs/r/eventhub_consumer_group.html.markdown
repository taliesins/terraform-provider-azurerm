---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_consumer_group"
sidebar_current: "docs-azurerm-resource-messaging-eventhub-consumer-group"
description: |-
  Manages a Event Hubs Consumer Group as a nested resource within an Event Hub.
---

# azurerm_eventhub_consumer_group

Manages a Event Hubs Consumer Group as a nested resource within an Event Hub.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acceptanceTestEventHubNamespace"
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
  capacity            = 2

  tags {
    environment = "Production"
  }
}

resource "azurerm_eventhub" "test" {
  name                = "acceptanceTestEventHub"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 2
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acceptanceTestEventHubConsumerGroup"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  eventhub_name       = "${azurerm_eventhub.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  user_metadata       = "some-meta-data"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Consumer Group resource. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the name of the grandparent EventHub Namespace. Changing this forces a new resource to be created.

* `eventhub_name` - (Required) Specifies the name of the EventHub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventHub Consumer Group's grandparent Namespace exists. Changing this forces a new resource to be created.

* `user_metadata` - (Optional) Specifies the user metadata.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub Consumer Group ID.

## Import

EventHub Consumer Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_consumer_group.consumerGroup1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumergroups/consumerGroup1
```
