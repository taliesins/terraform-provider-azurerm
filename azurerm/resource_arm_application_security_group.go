package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApplicationSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationSecurityGroupCreateUpdate,
		Read:   resourceArmApplicationSecurityGroupRead,
		Update: resourceArmApplicationSecurityGroupCreateUpdate,
		Delete: resourceArmApplicationSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmApplicationSecurityGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationSecurityGroupsClient

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	securityGroup := network.ApplicationSecurityGroup{
		Location: utils.String(location),
		Tags:     expandTags(tags),
	}
	_, createErr := client.CreateOrUpdate(resourceGroup, name, securityGroup, make(chan struct{}))
	err := <-createErr
	if err != nil {
		return fmt.Errorf("Error creating Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Application Security Group %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmApplicationSecurityGroupRead(d, meta)
}

func resourceArmApplicationSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationSecurityGroupsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["applicationSecurityGroups"]

	resp, err := client.Get(resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("resource_group_name", resourceGroup)
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmApplicationSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
