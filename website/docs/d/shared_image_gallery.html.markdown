---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image_gallery"
sidebar_current: "docs-azurerm-datasource-shared-image-gallery"
description: |-
  Gets information about an existing Shared Image Gallery.

---

# Data Source: azurerm_shared_image_gallery

Use this data source to access information about an existing Shared Image Gallery.

-> **NOTE** Shared Image Galleries are currently in Public Preview. You can find more information, including [how to register for the Public Preview here](https://azure.microsoft.com/en-gb/blog/announcing-the-public-preview-of-shared-image-gallery/).

## Example Usage

```hcl
data "azurerm_shared_image_gallery" "test" {
  name                = "my-image-gallery"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Shared Image Gallery.

* `resource_group_name` - (Required) The name of the Resource Group in which the Shared Image Gallery exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the Shared Image Gallery.

* `description` - A description for the Shared Image Gallery.

* `unique_name` - The unique name assigned to the Shared Image Gallery.

* `tags` - A mapping of tags which are assigned to the Shared Image Gallery.
