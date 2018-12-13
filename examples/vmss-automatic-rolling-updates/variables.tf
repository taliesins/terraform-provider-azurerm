# variable "subscription_id" {}
# variable "client_id" {}
# variable "client_secret" {}
# variable "tenant_id" {}

variable "resource_group_name" {
  description = "The name of the resource group in which to create the virtual network."
  default     = "tfex-vmss-ubuntu"
}

variable "location" {
  description = "The location/region where the virtual network is created. Changing this forces a new resource to be created."
  default     = "southcentralus"
}

variable "storage_account_tier" {
  description = "Defines the Tier of storage account to be created. Valid options are Standard and Premium."
  default     = "Standard"
}

variable "storage_replication_type" {
  description = "Defines the Replication Type to use for this storage account. Valid options include LRS, GRS etc."
  default     = "LRS"
}

variable "hostname" {
  description = "A string that determines the hostname/IP address of the origin server. This string could be a domain name, IPv4 address or IPv6 address."
  default     = "tfex-vmss-ubuntu"
}

variable "vm_sku" {
  description = "Size of VMs in the VM Scale Set."
  default     = "Standard_A1"
}

variable "ubuntu_os_version" {
  description = "The Ubuntu version for the VM. This will pick a fully patched image of this given Ubuntu version. Allowed values are: 15.10, 14.04.4-LTS."
  default     = "16.04.0-LTS"
}

variable "image_publisher" {
  description = "The name of the publisher of the image (az vm image list)"
  default     = "Canonical"
}

variable "image_offer" {
  description = "The name of the offer (az vm image list)"
  default     = "UbuntuServer"
}

variable "vmss_name_prefix" {
  description = "String used as a base for naming resources. Must be 1-9 characters in length for windows and 1-58 for linux images and globally unique across Azure. A hash is prepended to this string for some resources, and resource-specific information is appended."
}

variable "instance_count" {
  description = "Number of VM instances (100 or less)."
  default     = "5"
}

variable "admin_username" {
  description = "Admin username on all VMs."
}

variable "admin_password" {
  description = "Admin password on all VMs."
}
