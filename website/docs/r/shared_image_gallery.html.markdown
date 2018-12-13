---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_shared_image_gallery"
sidebar_current: "docs-azurerm-resource-compute-shared-image-gallery"
description: |-
  Manages a Shared Image Gallery.

---

# azurerm_shared_image_gallery

Manages a Shared Image Gallery.

-> **NOTE** Shared Image Galleries are currently in Public Preview. You can find more information, including [how to register for the Public Preview here](https://azure.microsoft.com/en-gb/blog/announcing-the-public-preview-of-shared-image-gallery/).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "example_image_gallery"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  description         = "Shared images and things."

  tags {
    Hello = "There"
    World = "Example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Shared Image Gallery. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Shared Image Gallery. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `description` - (Optional) A description for this Shared Image Gallery.

* `tags` - (Optional) A mapping of tags to assign to the Shared Image Gallery.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Shared Image Gallery.

## Import

Shared Image Galleries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_shared_image_gallery.gallery1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1
```
