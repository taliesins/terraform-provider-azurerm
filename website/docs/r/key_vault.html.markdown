---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault"
sidebar_current: "docs-azurerm-resource-key-vault-x"
description: |-
  Manages a Key Vault.
---

# azurerm_key_vault

Manages a Key Vault.

~> **NOTE:** It's possible to define Key Vault Access Policies both within [the `azurerm_key_vault` resource](key_vault.html) via the `access_policy` block and by using [the `azurerm_key_vault_access_policy` resource](key_vault_access_policy.html). However it's not possible to use both methods to manage Access Policies within a KeyVault, since there'll be conflicts.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_key_vault" "test" {
  name                        = "testvault"
  location                    = "${azurerm_resource_group.test.location}"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  enabled_for_disk_encryption = true
  tenant_id                   = "d6e396d0-5584-41dc-9fc0-268df99bc610"

  sku {
    name = "standard"
  }

  access_policy {
    tenant_id = "d6e396d0-5584-41dc-9fc0-268df99bc610"
    object_id = "d746815a-0433-4a21-b95d-fc437d2d475b"

    key_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
    ]
  }

  network_acls {
    default_action             = "Deny"
    bypass                     = "AzureServices"
  }

  tags {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Key Vault. Changing this forces a new resource to be created.

* `sku` - (Required) An SKU block as described below.

* `tenant_id` - (Required) The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault.

* `access_policy` - (Optional) An access policy block as described below. A maximum of 16 may be declared.
    
~> **NOTE:** It's possible to define Key Vault Access Policies both within [the `azurerm_key_vault` resource](key_vault.html) via the `access_policy` block and by using [the `azurerm_key_vault_access_policy` resource](key_vault_access_policy.html). However it's not possible to use both methods to manage Access Policies within a KeyVault, since there'll be conflicts.

* `enabled_for_deployment` - (Optional) Boolean flag to specify whether Azure Virtual Machines are permitted to retrieve certificates stored as secrets from the key vault. Defaults to `false`.

* `enabled_for_disk_encryption` - (Optional) Boolean flag to specify whether Azure Disk Encryption is permitted to retrieve secrets from the vault and unwrap keys. Defaults to `false`.

* `enabled_for_template_deployment` - (Optional) Boolean flag to specify whether Azure Resource Manager is permitted to retrieve secrets from the key vault. Defaults to `false`.

* `network_acls` - (Optional) A `network_acls` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `access_policy` block supports the following:

* `tenant_id` - (Required) The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault. Must match the `tenant_id` used above.

* `object_id` - (Required) The object ID of a user, service principal or security group in the Azure Active Directory tenant for the vault. The object ID must be unique for the list of access policies.

* `application_id` - (Optional) The object ID of an Application in Azure Active Directory.

* `certificate_permissions` - (Optional) List of certificate permissions, must be one or more from the following: `create`, `delete`, `deleteissuers`, `get`, `getissuers`, `import`, `list`, `listissuers`, `managecontacts`, `manageissuers`, `purge`, `recover`, `setissuers` and `update`.

* `key_permissions` - (Required) List of key permissions, must be one or more from the following: `backup`, `create`, `decrypt`, `delete`, `encrypt`, `get`, `import`, `list`, `purge`, `recover`, `restore`, `sign`, `unwrapKey`, `update`, `verify` and `wrapKey`.

* `secret_permissions` - (Required) List of secret permissions, must be one or more from the following: `backup`, `delete`, `get`, `list`, `purge`, `recover`, `restore` and `set`.


---

A `network_acls` block supports the following:

* `bypass` - (Required) Specifies which traffic can bypass the network rules. Possible values are `AzureServices` and `None`.

* `default_action` - (Required) The Default Action to use when no rules match from `ip_rules` / `virtual_network_subnet_ids`. Possible values are `Allow` and `Deny`.

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access thie Key Vault.

* `virtual_network_subnet_ids` - (Optional) One or more Subnet ID's which should be able to access this Key Vault.

---

A `sku` block supports the following:

* `name` - (Required) The Name of the SKU used for this Key Vault. Possible values are `Standard` and `Premium`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Key Vault.

* `vault_uri` - The URI of the Key Vault, used for performing operations on keys and secrets.

## Import

Key Vault's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/vault1
```
