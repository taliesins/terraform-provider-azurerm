data "azurerm_resource_group" "rg" {
  name = "${var.resource_group_name}"
}

resource "azurerm_virtual_network" "vnet" {
  name                = "vnet"
  location            = "${data.azurerm_resource_group.rg.location}"
  address_space       = ["10.0.0.0/16"]
  resource_group_name = "${data.azurerm_resource_group.rg.name}"
}

resource "azurerm_subnet" "subnet" {
  name                 = "${var.prefix}_subnet"
  virtual_network_name = "${azurerm_virtual_network.vnet.name}"
  resource_group_name  = "${data.azurerm_resource_group.rg.name}"
  address_prefix       = "10.0.10.0/24"
}

resource "azurerm_network_interface" "nic" {
  name                = "${var.prefix}_nic"
  location            = "${data.azurerm_resource_group.rg.location}"
  resource_group_name = "${data.azurerm_resource_group.rg.name}"

  ip_configuration {
    name                          = "${var.prefix}ipconfig"
    subnet_id                     = "${azurerm_subnet.subnet.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.pip.id}"
  }
}

resource "azurerm_public_ip" "pip" {
  name                         = "${var.prefix}-ip"
  location                     = "${data.azurerm_resource_group.rg.location}"
  resource_group_name          = "${data.azurerm_resource_group.rg.name}"
  public_ip_address_allocation = "Dynamic"
  domain_name_label            = "${var.dns_name}"
}

resource "azurerm_storage_account" "stor" {
  name                     = "${var.prefix}stor"
  location                 = "${data.azurerm_resource_group.rg.location}"
  resource_group_name      = "${data.azurerm_resource_group.rg.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_managed_disk" "datadisk" {
  name                 = "${var.hostname}-datadisk"
  location             = "${data.azurerm_resource_group.rg.location}"
  resource_group_name  = "${data.azurerm_resource_group.rg.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1023"
}

resource "azurerm_virtual_machine" "vm" {
  name                  = "${var.prefix}vm"
  location              = "${data.azurerm_resource_group.rg.location}"
  resource_group_name   = "${data.azurerm_resource_group.rg.name}"
  vm_size               = "${var.vm_size}"
  network_interface_ids = ["${azurerm_network_interface.nic.id}"]

  storage_image_reference {
    publisher = "${var.image_publisher}"
    offer     = "${var.image_offer}"
    sku       = "${var.image_sku}"
    version   = "${var.image_version}"
  }

  storage_os_disk {
    name              = "${var.hostname}-osdisk"
    managed_disk_type = "Standard_LRS"
    caching           = "ReadWrite"
    create_option     = "FromImage"
  }

  storage_data_disk {
    name              = "${var.hostname}-datadisk"
    managed_disk_id   = "${azurerm_managed_disk.datadisk.id}"
    managed_disk_type = "Standard_LRS"
    disk_size_gb      = "1023"
    create_option     = "Attach"
    lun               = 0
  }

  os_profile {
    computer_name  = "${var.hostname}"
    admin_username = "${var.admin_username}"
    admin_password = "${var.admin_password}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = "${azurerm_storage_account.stor.primary_blob_endpoint}"
  }
}
