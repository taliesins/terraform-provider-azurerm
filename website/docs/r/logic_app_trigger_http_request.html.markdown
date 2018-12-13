---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_trigger_http_request"
sidebar_current: "docs-azurerm-resource-logic-app-trigger-http-request"
description: |-
  Manages a HTTP Request Trigger within a Logic App Workflow
---

# azurerm_logic_app_trigger_http_request

Manages a HTTP Request Trigger within a Logic App Workflow

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "workflow-resources"
  location = "East US"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "workflow1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_logic_app_trigger_http_request" "test" {
  name         = "some-http-trigger"
  logic_app_id = "${azurerm_logic_app_workflow.test.id}"

  schema = <<SCHEMA
{
    "type": "object",
    "properties": {
        "hello": {
            "type": "string"
        }
    }
}
SCHEMA
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the HTTP Request Trigger to be created within the Logic App Workflow. Changing this forces a new resource to be created.

-> **NOTE:** This name must be unique across all Triggers within the Logic App Workflow.

* `logic_app_id` - (Required) Specifies the ID of the Logic App Workflow. Changing this forces a new resource to be created.

* `schema` - (Required) A JSON Blob defining the Schema of the incoming request. This needs to be valid JSON.

-> **NOTE:** To make the Trigger more readable, you may wish to consider using HEREDOC syntax (as shown above) or [the `local_file` resource](https://www.terraform.io/docs/providers/local/d/file.html) to load the schema from a file on disk.

* `method` - (Optional) Specifies the HTTP Method which the request be using. Possible values include `DELETE`, `GET`, `PATCH`, `POST` or `PUT`.

* `relative_path` - (Optional) Specifies the Relative Path used for this Request.

-> **NOTE:** When `relative_path` is set a `method` must also be set.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HTTP Request Trigger within the Logic App Workflow.

## Import

Logic App HTTP Request Triggers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_trigger_http_request.request1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Logic/workflows/workflow1/triggers/request1
```

-> **NOTE:** This ID is unique to Terraform and doesn't directly match to any other resource. To compose this ID, you can take the ID Logic App Workflow and append `/triggers/{name of the trigger}`.
