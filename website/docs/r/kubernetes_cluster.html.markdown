---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster"
sidebar_current: "docs-azurerm-resource-container-kubernetes-cluster"
description: |-
  Manages a managed Kubernetes Cluster (also known as AKS / Azure Kubernetes Service)
---

# azurerm_kubernetes_cluster

Manages a Managed Kubernetes Cluster (also known as AKS / Azure Kubernetes Service)

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

This example provisions a basic Managed Kubernetes Cluster. Other examples of the `azurerm_kubernetes_cluster` resource can be found in [the `./examples/kubernetes` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/kubernetes)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG1"
  location = "East US"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestagent1"

  agent_pool_profile {
    name            = "default"
    count           = 1
    vm_size         = "Standard_D1_v2"
    os_type         = "Linux"
    os_disk_size_gb = 30
  }

  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000000000000000000000000000"
  }

  tags {
    Environment = "Production"
  }
}

output "client_certificate" {
  value = "${azurerm_kubernetes_cluster.test.kube_config.0.client_certificate}"
}

output "kube_config" {
  value = "${azurerm_kubernetes_cluster.test.kube_config_raw}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Kubernetes Cluster to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Managed Kubernetes Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Managed Kubernetes Cluster should exist. Changing this forces a new resource to be created.

* `agent_pool_profile` - (Required) One or more `agent_pool_profile` blocks as documented below.

* `dns_prefix` - (Required) DNS prefix specified when creating the managed cluster.

* `service_principal` - (Required) A `service_principal` block as documented below.

---

* `addon_profile` - (Optional) A `addon_profile` block.

* `kubernetes_version` - (Optional) Version of Kubernetes specified when creating the AKS managed cluster. If not specified, the latest recommended version will be used at provisioning time (but won't auto-upgrade).

* `linux_profile` - (Optional) A `linux_profile` block.

* `network_profile` - (Optional) A `network_profile` block.

-> **NOTE:** If `network_profile` is not defined, `kubenet` profile will be used by default.

* `role_based_access_control` - (Optional) A `role_based_access_control` block. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `addon_profile` block supports the following:

* `http_application_routing` - (Optional) A `http_application_routing` block.
* `oms_agent` - (Optional) A `oms_agent` block. For more details, please visit [How to onboard Azure Monitor for containers](https://docs.microsoft.com/en-us/azure/monitoring/monitoring-container-insights-onboard).

---

A `agent_pool_profile` block supports the following:

* `name` - (Required) Unique name of the Agent Pool Profile in the context of the Subscription and Resource Group. Changing this forces a new resource to be created.
* `count` - (Required) Number of Agents (VMs) in the Pool. Possible values must be in the range of 1 to 100 (inclusive). Defaults to `1`.
* `vm_size` - (Required) The size of each VM in the Agent Pool (e.g. `Standard_F1`). Changing this forces a new resource to be created.

* `max_pods` - (Optional) The maximum number of pods that can run on each agent.
* `os_disk_size_gb` - (Optional) The Agent Operating System disk size in GB. Changing this forces a new resource to be created.
* `os_type` - (Optional) The Operating System used for the Agents. Possible values are `Linux` and `Windows`.  Changing this forces a new resource to be created. Defaults to `Linux`.
* `vnet_subnet_id` - (Optional) The ID of the Subnet where the Agents in the Pool should be provisioned. Changing this forces a new resource to be created.

~> **NOTE:** A route table should be configured on this Subnet.

---

A `azure_active_directory` block supports the following:

* `client_app_id` - (Required) The Client ID of an Azure Active Directory Application. Changing this forces a new resource to be created.

* `server_app_id` - (Required) The Server ID of an Azure Active Directory Application. Changing this forces a new resource to be created.

* `server_app_secret` - (Required) The Client Secret of an Azure Active Directory Application. Changing this forces a new resource to be created.

* `tenant_id` - (Optional) The Tenant ID used for Azure Active Directory Application. If this isn't specified the Tenant ID of the current Subscription is used. Changing this forces a new resource to be created.

---

A `http_application_routing` block supports the following:

* `enabled` (Required) Is HTTP Application Routing Enabled? Changing this forces a new resource to be created.

---

A `linux_profile` block supports the following:

* `admin_username` - (Required) The Admin Username for the Cluster. Changing this forces a new resource to be created.

* `ssh_key` - (Required) One or more `ssh_key` blocks. Changing this forces a new resource to be created.

---

A `network_profile` block supports the following:

* `network_plugin` - (Required) Network plugin to use for networking. Currently supported values are `azure` and `kubenet`. Changing this forces a new resource to be created.

-> **NOTE:** When `network_plugin` is set to `azure` - the `vnet_subnet_id` field in the `agent_pool_profile` block must be set.

* `dns_service_ip` - (Optional) IP address within the Kubernetes service address range that will be used by cluster service discovery (kube-dns). This is required when `network_plugin` is set to `kubenet`. Changing this forces a new resource to be created.

* `docker_bridge_cidr` - (Optional) IP address (in CIDR notation) used as the Docker bridge IP address on nodes. This is required when `network_plugin` is set to `kubenet`. Changing this forces a new resource to be created.

* `pod_cidr` - (Optional) The CIDR to use for pod IP addresses. This field can only be set when `network_plugin` is set to `kubenet`. Changing this forces a new resource to be created.

* `service_cidr` - (Optional) The Network Range used by the Kubernetes service. This is required when `network_plugin` is set to `kubenet`. Changing this forces a new resource to be created.

~> **NOTE:** This range should not be used by any network element on or connected to this VNet. Service address CIDR must be smaller than /12.

Examples of how to use [AKS with Advanced Networking](https://docs.microsoft.com/en-us/azure/aks/networking-overview#advanced-networking) can be [found in the `./examples/kubernetes/` directory in the Github repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/kubernetes).

---

A `oms_agent` block supports the following:

* `enabled` - (Required) Is the OMS Agent Enabled?

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace which the OMS Agent should send data to.

---

A `role_based_access_control` block supports the following:

* `azure_active_directory` - (Required) An `azure_active_directory` block. Changing this forces a new resource to be created.

* `enabled` - (Required) Is Role Based Access Control Enabled? Changing this forces a new resource to be created.

---

A `service_principal` block supports the following:

* `client_id` - (Required) The Client ID for the Service Principal. Changing this forces a new resource to be created.

* `client_secret` - (Required) The Client Secret for the Service Principal. Changing this forces a new resource to be created.

---

A `ssh_key` block supports the following:

* `key_data` - (Required) The Public SSH Key used to access the cluster. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The Kubernetes Managed Cluster ID.

* `fqdn` - The FQDN of the Azure Kubernetes Managed Cluster.

* `kube_admin_config` - A `kube_admin_config` block as defined below. This is only available when Role Based Access Control with Azure Active Directory is enabled.

* `kube_admin_config_raw` - Raw Kubernetes config for the admin account to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools. This is only available when Role Based Access Control with Azure Active Directory is enabled.

* `kube_config` - A `kube_config` block as defined below.

* `kube_config_raw` - Raw Kubernetes config to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools

* `http_application_routing` - A `http_application_routing` block as defined below.

* `node_resource_group` - The auto-generated Resource Group which contains the resources for this Managed Kubernetes Cluster.

---

A `http_application_routing` block exports the following:

* `http_application_routing_zone_name` - The Zone Name of the HTTP Application Routing.

---

The `kube_admin_config` and `kube_config` blocks export the following::

* `client_key` - Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.

* `client_certificate` - Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.

* `cluster_ca_certificate` - Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.

* `host` - The Kubernetes cluster server host.

* `username` - A username used to authenticate to the Kubernetes cluster.

* `password` - A password or token used to authenticate to the Kubernetes cluster.

-> **NOTE:** It's possible to use these credentials with [the Kubernetes Provider](/docs/providers/kubernetes/index.html) like so:

```
provider "kubernetes" {
  host                   = "${azurerm_kubernetes_cluster.main.kube_config.0.host}"
  username               = "${azurerm_kubernetes_cluster.main.kube_config.0.username}"
  password               = "${azurerm_kubernetes_cluster.main.kube_config.0.password}"
  client_certificate     = "${base64decode(azurerm_kubernetes_cluster.main.kube_config.0.client_certificate)}"
  client_key             = "${base64decode(azurerm_kubernetes_cluster.main.kube_config.0.client_key)}"
  cluster_ca_certificate = "${base64decode(azurerm_kubernetes_cluster.main.kube_config.0.cluster_ca_certificate)}"
}
```

---

## Import

Managed Kubernetes Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1
```
