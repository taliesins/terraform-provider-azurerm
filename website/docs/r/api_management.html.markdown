---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management"
sidebar_current: "docs-azurerm-resource-api-management-x"
description: |-
  Manages an API Management Service.
---

# azurerm_api_management

Manages an API Management Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "test" {
  name                = "example-apim"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku {
    name     = "Developer"
    capacity = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the API Management Service exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service should be exist. Changing this forces a new resource to be created.

* `publisher_name` - (Required) The name of publisher/company.

* `publisher_email` - (Required) The email of publisher/company.

* `sku` - (Required) A `sku` block as documented below.

---

* `additional_location` - (Optional) One or more `additional_location` blocks as defined below.

* `certificate` - (Optional) One or more (up to 10) `certificate` blocks as defined below.

* `identity` - (Optional) An `identity` block is documented below.

* `hostname_configuration` - (Optional) A `hostname_configuration` block as defined below.

* `notification_sender_email` - (Optional) Email address from which the notification will be sent.

* `security` - (Optional) A `security` block as defined below.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

A `additional_location` block supports the following:

* `location` - (Required) The name of the Azure Region in which the API Management Service should be expanded to.

---

A `certificate` block supports the following:

* `encoded_certificate` - (Required) The Base64 Encoded PFX Certificate.

* `certificate_password` - (Required) The password for the certificate.

* `store_name` - (Required) The name of the Certificate Store where this certificate should be stored. Possible values are `CertificateAuthority` and `Root`.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this API Management Service. At this time the only supported value is`SystemAssigned`.

---

A `security` block supports the following:

* `disable_backend_ssl30` - (Optional) Should SSL 3.0 be disabled on the backend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30` field

* `disable_backend_tls10` - (Optional) Should TLS 1.0 be disabled on the backend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10` field

* `disable_backend_tls11` - (Optional) Should TLS 1.1 be disabled on the backend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11` field

* `disable_frontend_ssl30` - (Optional) Should SSL 3.0 be disabled on the frontend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30` field

* `disable_frontend_tls10` - (Optional) Should TLS 1.0 be disabled on the frontend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10` field

* `disable_frontend_tls11` - (Optional) Should TLS 1.1 be disabled on the frontend of the gateway? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11` field

* `disable_triple_des_chipers` - (Optional) Should the `TLS_RSA_WITH_3DES_EDE_CBC_SHA` cipher be disabled for alL TLS versions (1.0, 1.1 and 1.2)? Defaults to `false`.

-> **info:** This maps to the `Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168` field

---

A `sku` block supports the following:

* `name` - (Required) Specifies the Pricing Tier for the API Management Service. Possible values include: `Developer`, `Basic`, `Standard` and `Premium`.

* `capacity` - (Required) Specifies the Pricing Capacity for the API Management Service.

---

A `hostname_configuration` block supports the following:

* `management` - (Optional) One or more `management` blocks as documented below.

* `portal` - (Optional) One or more `portal` blocks as documented below.

* `proxy` - (Optional) One or more `proxy` blocks as documented below.

* `scm` - (Optional) One or more `scm` blocks as documented below.

---

A `management`, `portal` and `scm` block supports the following:

* `host_name` - (Required) The Hostname to use for the Management API.

* `key_vault_id` - (Optional) The ID of the Key Vault Secret containing the SSL Certificate, which must be should be of the type `application/x-pkcs12`.

-> **NOTE:** Setting this field requires the `identity` block to be specified, since this identity is used for to retrieve the Key Vault Certificate. Auto-updating the Certificate from the Key Vault requires the Secret version isn't specified.

* `certificate` - (Optional) The Base64 Encoded Certificate.

* `certificate_password` - (Optional) The password associated with the certificate provided above.

-> **NOTE:** Either `key_vault_id` or `certificate` and `certificate_password` must be specified.

* `negotiate_client_certificate` - (Optional) Should Client Certificate Negotiation be enabled for this Hostname? Defaults to `false`.

---

A `proxy` block supports the following:

* `default_ssl_binding` - (Optional) Is the certificate associated with this Hostname the Default SSL Certificate? This is used when an SNI header isn't specified by a client. Defaults to `false`.

* `host_name` - (Required) The Hostname to use for the Management API.

* `key_vault_id` - (Optional) The ID of the Key Vault Secret containing the SSL Certificate, which must be should be of the type `application/x-pkcs12`.

-> **NOTE:** Setting this field requires the `identity` block to be specified, since this identity is used for to retrieve the Key Vault Certificate. Auto-updating the Certificate from the Key Vault requires the Secret version isn't specified.

* `certificate` - (Optional) The Base64 Encoded Certificate.

* `certificate_password` - (Optional) The password associated with the certificate provided above.

-> **NOTE:** Either `key_vault_id` or `certificate` and `certificate_password` must be specified.

* `negotiate_client_certificate` - (Optional) Should Client Certificate Negotiation be enabled for this Hostname? Defaults to `false`.



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Service.

* `gateway_url` - The URL of the Gateway for the API Management Service.

* `gateway_regional_url` - The Region URL for the Gateway of the API Management Service.

* `management_api_url` - The URL for the Management API associated with this API Management service.

* `portal_url` - The URL for the Publisher Portal associated with this API Management service.

* `public_ip_addresses` - The Public IP addresses of the API Management Service.

* `scm_url` - The URL for the SCM (Source Code Management) Endpoint associated with this API Management service.

* `identity` - An `identity` block as defined below.

* `additional_location` - One or more `additional_location` blocks as documented below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

An `additional_location` block exports the following:

* `gateway_regional_url` - The URL of the Regional Gateway for the API Management Service in the specified region.

* `public_ip_addresses` - Public Static Load Balanced IP addresses of the API Management service in the additional location. Available only for Basic, Standard and Premium SKU.

## Import

API Management Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1
```
