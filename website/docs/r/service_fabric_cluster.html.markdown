---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_fabric_cluster"
sidebar_current: "docs-azurerm-resource-service-fabric-cluster"
description: |-
  Manage a Service Fabric Cluster.
---

# azurerm_service_fabric_cluster

Manage a Service Fabric Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_fabric_cluster" "test" {
  name                 = "example-servicefabric"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  reliability_level    = "Bronze"
  upgrade_mode         = "Manual"
  cluster_code_version = "6.3.176.9494"
  vm_image             = "Windows"
  management_endpoint  = "https://example:80"

  node_type {
    name                 = "first"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Fabric Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Service Fabric Cluster exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Service Fabric Cluster should exist. Changing this forces a new resource to be created.

* `reliability_level` - (Required) Specifies the Reliability Level of the Cluster. Possible values include `None`, `Bronze`, `Silver`, `Gold` and `Platinum`.

-> **NOTE:** The Reliability Level of the Cluster depends on the number of nodes in the Cluster: `Platinum` requires at least 9 VM's, `Gold` requires at least 7 VM's, `Silver` requires at least 5 VM's, `Bronze` requires at least 3 VM's.

* `management_endpoint` - (Required) Specifies the Management Endpoint of the cluster such as `http://example.com`. Changing this forces a new resource to be created.

* `node_type` - (Required) One or more `node_type` blocks as defined below.

* `upgrade_mode` - (Required) Specifies the Upgrade Mode of the cluster. Possible values are `Automatic` or `Manual`.

* `vm_image` - (Required) Specifies the Image expected for the Service Fabric Cluster, such as `Windows`. Changing this forces a new resource to be created.

---

* `cluster_code_version` - (Optional) Required if Upgrade Mode set to `Manual`, Specifies the Version of the Cluster Code of the cluster.

* `add_on_features` - (Optional) A List of one or more features which should be enabled, such as `DnsService`.

* `certificate` - (Optional) A `certificate` block as defined below.

* `client_certificate_thumbprint` - (Optional) One or two `client_certificate_thumbprint` blocks as defined below.

-> **NOTE:** If Client Certificates are enabled then at a Certificate must be configured on the cluster.

* `diagnostics_config` - (Optional) A `diagnostics_config` block as defined below. Changing this forces a new resource to be created.

* `fabric_settings` - (Optional) One or more `fabric_settings` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `certificate` block supports the following:

* `thumbprint` - (Required) The Thumbprint of the Certificate.

* `thumbprint_secondary` - (Required) The Secondary Thumbprint of the Certificate.

* `x509_store_name` - (Required) The X509 Store where the Certificate Exists, such as `My`.

---

A `client_certificate_thumbprint` block supports the following:

* `thumbprint` - (Required) The Thumbprint associated with the Client Certificate.

* `is_admin` - (Required) Does the Client Certificate have Admin Access to the cluster? Non-admin clients can only perform read only operations on the cluster.

---

A `diagnostics_config` block supports the following:

* `storage_account_name` - (Required) The name of the Storage Account where the Diagnostics should be sent to.

* `protected_account_key_name` - (Required) The protected diagnostics storage key name, such as `StorageAccountKey1`.

* `blob_endpoint` - (Required) The Blob Endpoint of the Storage Account.

* `queue_endpoint` - (Required) The Queue Endpoint of the Storage Account.

* `table_endpoint` - (Required) The Table Endpoint of the Storage Account.

---

A `fabric_settings` block supports the following:

* `name` - (Required) The name of the Fabric Setting, such as `Security` or `Federation`.

* `parameters` - (Optional) A map containing settings for the specified Fabric Setting.

---

A `node_type` block supports the following:

* `name` - (Required) The name of the Node Type. Changing this forces a new resource to be created.

* `instance_count` - (Required) The number of nodes for this Node Type.

* `is_primary` - (Required) Is this the Primary Node Type? Changing this forces a new resource to be created.

* `client_endpoint_port` - (Required) The Port used for the Client Endpoint for this Node Type. Changing this forces a new resource to be created.

* `http_endpoint_port` - (Required) The Port used for the HTTP Endpoint for this Node Type. Changing this forces a new resource to be created.

* `durability_level` - (Optional) The Durability Level for this Node Type. Possible values include `Bronze`, `Gold` and `Silver`. Defaults to `Bronze`. Changing this forces a new resource to be created.

* `application_ports` - (Optional) A `application_ports` block as defined below.

* `ephemeral_ports` - (Optional) A `ephemeral_ports` block as defined below.

---

A `application_ports` block supports the following:

* `start_port` - (Required) The start of the Application Port Range on this Node Type.

* `end_port` - (Required) The end of the Application Port Range on this Node Type.

---

A `ephemeral_ports` block supports the following:

* `start_port` - (Required) The start of the Ephemeral Port Range on this Node Type.

* `end_port` - (Required) The end of the Ephemeral Port Range on this Node Type.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Service Fabric Cluster.

* `cluster_endpoint` - The Cluster Endpoint for this Service Fabric Cluster.

## Import

Service Fabric Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_fabric_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceFabric/clusters/cluster1
```
