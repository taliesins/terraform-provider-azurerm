---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_store_file"
sidebar_current: "docs-azurerm-resource-data-lake-store-file"
description: |-
  Manage a Azure Data Lake Store File.
---

# azurerm_data_lake_store_file

Manage a Azure Data Lake Store File.

~> **Note:** If you want to change the data in the remote file without changing the `local_file_path`, then 
taint the resource so the `azurerm_data_lake_store_file` gets recreated with the new data.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "northeurope"
}

resource "azurerm_data_lake_store" "example" {
  name                = "consumptiondatalake"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}

resource "azurerm_data_lake_store_file" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  local_file_path     = "/path/to/local/file"
  remote_file_path    = "/path/created/for/remote/file"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) Specifies the name of the Data Lake Store for which the File should created.

* `local_file_path` - (Required) The path to the local file to be added to the Data Lake Store.

* `remote_file_path` - (Required) The path created for the file on the Data Lake Store.

## Import

Date Lake Store File's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_lake_store_file.test example.azuredatalakestore.net/test/example.txt
```
