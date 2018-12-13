---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub"
sidebar_current: "docs-azurerm-resource-messaging-iothub-x"
description: |-
  Manages an IotHub
---

# azurerm_iothub

Manages an IotHub

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_storage_account" "test" {
  name                     = "teststa"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "test"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  endpoint {
    type                       = "AzureIotHub.StorageContainer"
    connection_string          = "${azurerm_storage_account.test.primary_blob_connection_string}"
    name                       = "export"
    batch_frequency_in_seconds = 60
    max_chunk_size_in_bytes    = 10485760
    container_name             = "test"
    encoding                   = "Avro"
    file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
  }

  route {
    name           = "export"
    source         = "DeviceMessages"
    condition      = "true"
    endpoint_names = ["export"]
    enabled        = true
  }

  tags {
    "purpose" = "testing"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the IotHub resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the IotHub resource has to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource has to be createc. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `endpoint` - (Optional) An `endpoint` block as defined below.

* `route` - (Optional) A `route` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) The name of the sku. Possible values are `B1`, `B2`, `B3`, `F1`, `S1`, `S2`, and `S3`.

* `tier` - (Required) The billing tier for the IoT Hub. Possible values are `Basic`, `Free` or `Standard`.

~> **NOTE:** Only one IotHub can be on the `Free` tier per subscription.

* `capacity` - (Required) The number of provisioned IoT Hub units.

---

An `endpoint` block supports the following:

* `type` - (Required) The type of the endpoint. Possible values are `AzureIotHub.StorageContainer`, `AzureIotHub.ServiceBusQueue`, `AzureIotHub.ServiceBusTopic` or `AzureIotHub.EventHub`.

* `connection_string` - (Required) The connection string for the endpoint.

* `name` - (Required) The name of the endpoint. The name must be unique across endpoint types. The following names are reserved:  `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`.

* `batch_frequency_in_seconds` - (Optional) Time interval at which blobs are written to storage. Value should be between 60 and 720 seconds. Default value is 300 seconds. This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

* `max_chunk_size_in_bytes` - (Optional) Maximum number of bytes for each blob written to storage. Value should be between 10485760(10MB) and 524288000(500MB). Default value is 314572800(300MB). This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

* `container_name` - (Optional) The name of storage container in the storage account. This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

* `encoding` - (Optional) Encoding that is used to serialize messages to blobs. Supported values are 'avro' and 'avrodeflate'. Default value is 'avro'. This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

* `file_name_format` - (Optional) File name format for the blob. Default format is ``{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}``. All parameters are mandatory but can be reordered. This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

---

A `route` block supports the following:

* `name` - (Required) The name of the route. The name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.

* `source` - (Required) The source that the routing rule is to be applied to, such as `DeviceMessages`. Possible values include: `RoutingSourceInvalid`, `RoutingSourceDeviceMessages`, `RoutingSourceTwinChangeEvents`, `RoutingSourceDeviceLifecycleEvents`, `RoutingSourceDeviceJobLifecycleEvents`.

* `condition` - (Optional) The condition that is evaluated to apply the routing rule. If no condition is provided, it evaluates to true by default. For grammar, see: https://docs.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language.

* `endpoint_names` - (Required) The list of endpoints to which messages that satisfy the condition are routed.

* `enabled` - (Required) Used to specify whether a route is enabled.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTHub.

* `event_hub_events_endpoint` -  The EventHub compatible endpoint for events data
* `event_hub_events_path` -  The EventHub compatible path for events data
* `event_hub_operations_endpoint` -  The EventHub compatible endpoint for operational data
* `event_hub_operations_path` -  The EventHub compatible path for operational data

-> **NOTE:** These fields can be used in conjunction with the `shared_access_policy` block to build a connection string

* `hostname` - The hostname of the IotHub Resource.

* `shared_access_policy` - One or more `shared_access_policy` blocks as defined below.

---

A `shared access policy` block contains the following:

* `key_name` - The name of the shared access policy.

* `primary_key` - The primary key.

* `secondary_key` - The secondary key.

* `permissions` - The permissions assigned to the shared access policy.

## Import

IoTHubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub.hub1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1
```
