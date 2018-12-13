package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerGroupCreate,
		Read:   resourceArmContainerGroupRead,
		Delete: resourceArmContainerGroupDelete,
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

			"ip_address_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "Public",
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.Public),
				}, true),
			},

			"os_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.Windows),
					string(containerinstance.Linux),
				}, true),
			},

			"image_registry_credential": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
							ForceNew:     true,
						},

						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
							ForceNew:     true,
						},

						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
							ForceNew:     true,
						},
					},
				},
			},

			"tags": tagsForceNewSchema(),

			"restart_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          string(containerinstance.Always),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.Always),
					string(containerinstance.Never),
					string(containerinstance.OnFailure),
				}, true),
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_name_label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"container": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"image": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"cpu": {
							Type:     schema.TypeFloat,
							Required: true,
							ForceNew: true,
						},

						"memory": {
							Type:     schema.TypeFloat,
							Required: true,
							ForceNew: true,
						},

						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"protocol": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerinstance.TCP),
								string(containerinstance.UDP),
							}, true),
						},

						"environment_variables": {
							Type:     schema.TypeMap,
							ForceNew: true,
							Optional: true,
						},

						"secure_environment_variables": {
							Type:      schema.TypeMap,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						"command": {
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Use `commands` instead.",
						},

						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"volume": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"mount_path": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
										Default:  false,
									},

									"share_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"storage_account_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"storage_account_key": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	containerGroupsClient := meta.(*ArmClient).containerGroupsClient

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	OSType := d.Get("os_type").(string)
	IPAddressType := d.Get("ip_address_type").(string)
	tags := d.Get("tags").(map[string]interface{})
	restartPolicy := d.Get("restart_policy").(string)

	containers, containerGroupPorts, containerGroupVolumes := expandContainerGroupContainers(d)
	containerGroup := containerinstance.ContainerGroup{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			Containers:    containers,
			RestartPolicy: containerinstance.ContainerGroupRestartPolicy(restartPolicy),
			IPAddress: &containerinstance.IPAddress{
				Type:  containerinstance.ContainerGroupIPAddressType(IPAddressType),
				Ports: containerGroupPorts,
			},
			OsType:                   containerinstance.OperatingSystemTypes(OSType),
			Volumes:                  containerGroupVolumes,
			ImageRegistryCredentials: expandContainerImageRegistryCredentials(d),
		},
	}

	if dnsNameLabel := d.Get("dns_name_label").(string); dnsNameLabel != "" {
		containerGroup.ContainerGroupProperties.IPAddress.DNSNameLabel = &dnsNameLabel
	}

	if _, err := containerGroupsClient.CreateOrUpdate(ctx, resGroup, name, containerGroup); err != nil {
		return err
	}

	read, err := containerGroupsClient.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read container group %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerGroupRead(d, meta)
}

func resourceArmContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).containerGroupsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["containerGroups"]

	resp, err := client.Get(ctx, resourceGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Group %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ContainerGroupProperties; props != nil {
		containerConfigs := flattenContainerGroupContainers(d, resp.Containers, props.IPAddress.Ports, props.Volumes)
		if err := d.Set("container", containerConfigs); err != nil {
			return fmt.Errorf("Error setting `container`: %+v", err)
		}

		if err := d.Set("image_registry_credential", flattenContainerImageRegistryCredentials(d, props.ImageRegistryCredentials)); err != nil {
			return fmt.Errorf("Error setting `capabilities`: %+v", err)
		}

		if address := props.IPAddress; address != nil {
			d.Set("ip_address_type", address.Type)
			d.Set("ip_address", address.IP)
			d.Set("dns_name_label", address.DNSNameLabel)
			d.Set("fqdn", address.Fqdn)
		}

		d.Set("restart_policy", string(props.RestartPolicy))
		d.Set("os_type", string(props.OsType))
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).containerGroupsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["containerGroups"]

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Container Group %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandContainerGroupContainers(d *schema.ResourceData) (*[]containerinstance.Container, *[]containerinstance.Port, *[]containerinstance.Volume) {
	containersConfig := d.Get("container").([]interface{})
	containers := make([]containerinstance.Container, 0)
	containerGroupPorts := make([]containerinstance.Port, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)

	for _, containerConfig := range containersConfig {
		data := containerConfig.(map[string]interface{})

		name := data["name"].(string)
		image := data["image"].(string)
		cpu := data["cpu"].(float64)
		memory := data["memory"].(float64)

		container := containerinstance.Container{
			Name: utils.String(name),
			ContainerProperties: &containerinstance.ContainerProperties{
				Image: utils.String(image),
				Resources: &containerinstance.ResourceRequirements{
					Requests: &containerinstance.ResourceRequests{
						MemoryInGB: utils.Float(memory),
						CPU:        utils.Float(cpu),
					},
				},
			},
		}

		if v := data["port"]; v != 0 {
			port := int32(v.(int))

			// container port (port number)
			container.Ports = &[]containerinstance.ContainerPort{
				{
					Port: &port,
				},
			}

			// container group port (port number + protocol)
			containerGroupPort := containerinstance.Port{
				Port: &port,
			}

			if v, ok := data["protocol"]; ok {
				protocol := v.(string)
				containerGroupPort.Protocol = containerinstance.ContainerGroupNetworkProtocol(strings.ToUpper(protocol))
			}

			containerGroupPorts = append(containerGroupPorts, containerGroupPort)
		}

		// Set both sensitive and non-secure environment variables
		var envVars *[]containerinstance.EnvironmentVariable
		var secEnvVars *[]containerinstance.EnvironmentVariable

		// Expand environment_variables into slice
		if v, ok := data["environment_variables"]; ok {
			envVars = expandContainerEnvironmentVariables(v, false)
		}

		// Expand secure_environment_variables into slice
		if v, ok := data["secure_environment_variables"]; ok {
			secEnvVars = expandContainerEnvironmentVariables(v, true)
		}

		// Combine environment variable slices
		*envVars = append(*envVars, *secEnvVars...)

		// Set both secure and non secure environment variables
		container.EnvironmentVariables = envVars

		if v, ok := data["commands"]; ok {
			c := v.([]interface{})
			command := make([]string, 0)
			for _, v := range c {
				command = append(command, v.(string))
			}

			container.Command = &command
		}

		if container.Command == nil {
			if v := data["command"]; v != "" {
				command := strings.Split(v.(string), " ")
				container.Command = &command
			}
		}

		if v, ok := data["volume"]; ok {
			volumeMounts, containerGroupVolumesPartial := expandContainerVolumes(v)
			container.VolumeMounts = volumeMounts
			if containerGroupVolumesPartial != nil {
				containerGroupVolumes = append(containerGroupVolumes, *containerGroupVolumesPartial...)
			}
		}

		containers = append(containers, container)
	}

	return &containers, &containerGroupPorts, &containerGroupVolumes
}

func expandContainerEnvironmentVariables(input interface{}, secure bool) *[]containerinstance.EnvironmentVariable {

	envVars := input.(map[string]interface{})
	output := make([]containerinstance.EnvironmentVariable, 0, len(envVars))

	if secure {

		for k, v := range envVars {
			ev := containerinstance.EnvironmentVariable{
				Name:        utils.String(k),
				SecureValue: utils.String(v.(string)),
			}

			output = append(output, ev)
		}

	} else {

		for k, v := range envVars {
			ev := containerinstance.EnvironmentVariable{
				Name:  utils.String(k),
				Value: utils.String(v.(string)),
			}

			output = append(output, ev)
		}
	}
	return &output
}

func expandContainerImageRegistryCredentials(d *schema.ResourceData) *[]containerinstance.ImageRegistryCredential {
	credsRaw := d.Get("image_registry_credential").([]interface{})
	if len(credsRaw) == 0 {
		return nil
	}

	output := make([]containerinstance.ImageRegistryCredential, 0, len(credsRaw))

	for _, c := range credsRaw {
		credConfig := c.(map[string]interface{})

		output = append(output, containerinstance.ImageRegistryCredential{
			Server:   utils.String(credConfig["server"].(string)),
			Password: utils.String(credConfig["password"].(string)),
			Username: utils.String(credConfig["username"].(string)),
		})
	}

	return &output
}

func expandContainerVolumes(input interface{}) (*[]containerinstance.VolumeMount, *[]containerinstance.Volume) {
	volumesRaw := input.([]interface{})

	if len(volumesRaw) == 0 {
		return nil, nil
	}

	volumeMounts := make([]containerinstance.VolumeMount, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)

	for _, volumeRaw := range volumesRaw {
		volumeConfig := volumeRaw.(map[string]interface{})

		name := volumeConfig["name"].(string)
		mountPath := volumeConfig["mount_path"].(string)
		readOnly := volumeConfig["read_only"].(bool)
		shareName := volumeConfig["share_name"].(string)
		storageAccountName := volumeConfig["storage_account_name"].(string)
		storageAccountKey := volumeConfig["storage_account_key"].(string)

		vm := containerinstance.VolumeMount{
			Name:      utils.String(name),
			MountPath: utils.String(mountPath),
			ReadOnly:  utils.Bool(readOnly),
		}

		volumeMounts = append(volumeMounts, vm)

		cv := containerinstance.Volume{
			Name: utils.String(name),
			AzureFile: &containerinstance.AzureFileVolume{
				ShareName:          utils.String(shareName),
				ReadOnly:           utils.Bool(readOnly),
				StorageAccountName: utils.String(storageAccountName),
				StorageAccountKey:  utils.String(storageAccountKey),
			},
		}

		containerGroupVolumes = append(containerGroupVolumes, cv)
	}

	return &volumeMounts, &containerGroupVolumes
}

func flattenContainerImageRegistryCredentials(d *schema.ResourceData, input *[]containerinstance.ImageRegistryCredential) []interface{} {
	if input == nil {
		return nil
	}
	configsOld := d.Get("image_registry_credential").([]interface{})

	output := make([]interface{}, 0)
	for i, cred := range *input {
		credConfig := make(map[string]interface{})
		if cred.Server != nil {
			credConfig["server"] = *cred.Server
		}
		if cred.Username != nil {
			credConfig["username"] = *cred.Username
		}

		if len(configsOld) > i {
			data := configsOld[i].(map[string]interface{})
			oldServer := data["server"].(string)
			if cred.Server != nil && *cred.Server == oldServer {
				if v, ok := d.GetOk(fmt.Sprintf("image_registry_credential.%d.password", i)); ok {
					credConfig["password"] = v.(string)
				}
			}
		}

		output = append(output, credConfig)
	}
	return output
}

func flattenContainerGroupContainers(d *schema.ResourceData, containers *[]containerinstance.Container, containerGroupPorts *[]containerinstance.Port, containerGroupVolumes *[]containerinstance.Volume) []interface{} {

	//map old container names to index so we can look up things up
	nameIndexMap := map[string]int{}
	for i, c := range d.Get("container").([]interface{}) {
		cfg := c.(map[string]interface{})
		nameIndexMap[cfg["name"].(string)] = i

	}

	containerCfg := make([]interface{}, 0, len(*containers))
	for _, container := range *containers {

		//TODO fix this crash point
		name := *container.Name

		//get index from name
		index := nameIndexMap[name]

		containerConfig := make(map[string]interface{})
		containerConfig["name"] = name

		if v := container.Image; v != nil {
			containerConfig["image"] = *v
		}

		if resources := container.Resources; resources != nil {
			if resourceRequests := resources.Requests; resourceRequests != nil {
				if v := resourceRequests.CPU; v != nil {
					containerConfig["cpu"] = *v
				}
				if v := resourceRequests.MemoryInGB; v != nil {
					containerConfig["memory"] = *v
				}
			}
		}

		if len(*container.Ports) > 0 {
			containerPort := *(*container.Ports)[0].Port
			containerConfig["port"] = containerPort
			// protocol isn't returned in container config, have to search in container group ports
			protocol := ""
			if containerGroupPorts != nil {
				for _, cgPort := range *containerGroupPorts {
					if *cgPort.Port == containerPort {
						protocol = string(cgPort.Protocol)
					}
				}
			}
			if protocol != "" {
				containerConfig["protocol"] = protocol
			}
		}

		if container.EnvironmentVariables != nil {
			if len(*container.EnvironmentVariables) > 0 {
				containerConfig["environment_variables"] = flattenContainerEnvironmentVariables(container.EnvironmentVariables, false, d, index)
			}
		}

		if container.EnvironmentVariables != nil {
			if len(*container.EnvironmentVariables) > 0 {
				containerConfig["secure_environment_variables"] = flattenContainerEnvironmentVariables(container.EnvironmentVariables, true, d, index)
			}
		}

		commands := make([]string, 0)
		if command := container.Command; command != nil {
			containerConfig["command"] = strings.Join(*command, " ")
			commands = *command
		}
		containerConfig["commands"] = commands

		if containerGroupVolumes != nil && container.VolumeMounts != nil {
			// Also pass in the container volume config from schema
			var containerVolumesConfig *[]interface{}
			containersConfigRaw := d.Get("container").([]interface{})
			for _, containerConfigRaw := range containersConfigRaw {
				data := containerConfigRaw.(map[string]interface{})
				nameRaw := data["name"].(string)
				if nameRaw == *container.Name {
					// found container config for current container
					// extract volume mounts from config
					if v, ok := data["volume"]; ok {
						containerVolumesRaw := v.([]interface{})
						containerVolumesConfig = &containerVolumesRaw
					}
				}
			}
			containerConfig["volume"] = flattenContainerVolumes(container.VolumeMounts, containerGroupVolumes, containerVolumesConfig)
		}

		containerCfg = append(containerCfg, containerConfig)
	}

	return containerCfg
}

func flattenContainerEnvironmentVariables(input *[]containerinstance.EnvironmentVariable, isSecure bool, d *schema.ResourceData, oldContainerIndex int) map[string]interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return output
	}

	if isSecure {
		for _, envVar := range *input {

			if envVar.Name != nil && envVar.Value == nil {
				if v, ok := d.GetOk(fmt.Sprintf("container.%d.secure_environment_variables.%s", oldContainerIndex, *envVar.Name)); ok {
					log.Printf("[DEBUG] SECURE    : Name: %s - Value: %s", *envVar.Name, v.(string))
					output[*envVar.Name] = v.(string)
				}
			}
		}
	} else {
		for _, envVar := range *input {
			if envVar.Name != nil && envVar.Value != nil {
				log.Printf("[DEBUG] NOT SECURE: Name: %s - Value: %s", *envVar.Name, *envVar.Value)
				output[*envVar.Name] = *envVar.Value
			}
		}
	}

	return output
}

func flattenContainerVolumes(volumeMounts *[]containerinstance.VolumeMount, containerGroupVolumes *[]containerinstance.Volume, containerVolumesConfig *[]interface{}) []interface{} {
	volumeConfigs := make([]interface{}, 0)

	if volumeMounts == nil {
		return volumeConfigs
	}

	for _, vm := range *volumeMounts {
		volumeConfig := make(map[string]interface{})
		if vm.Name != nil {
			volumeConfig["name"] = *vm.Name
		}
		if vm.MountPath != nil {
			volumeConfig["mount_path"] = *vm.MountPath
		}
		if vm.ReadOnly != nil {
			volumeConfig["read_only"] = *vm.ReadOnly
		}

		// find corresponding volume in container group volumes
		// and use the data
		if containerGroupVolumes != nil {
			for _, cgv := range *containerGroupVolumes {
				if cgv.Name == nil || vm.Name == nil {
					continue
				}

				if *cgv.Name == *vm.Name {
					if file := cgv.AzureFile; file != nil {
						if file.ShareName != nil {
							volumeConfig["share_name"] = *file.ShareName
						}
						if file.StorageAccountName != nil {
							volumeConfig["storage_account_name"] = *file.StorageAccountName
						}
						// skip storage_account_key, is always nil
					}
				}
			}
		}

		// find corresponding volume in config
		// and use the data
		if containerVolumesConfig != nil {
			for _, cvr := range *containerVolumesConfig {
				cv := cvr.(map[string]interface{})
				rawName := cv["name"].(string)
				if vm.Name != nil && *vm.Name == rawName {
					storageAccountKey := cv["storage_account_key"].(string)
					volumeConfig["storage_account_key"] = storageAccountKey
				}
			}
		}

		volumeConfigs = append(volumeConfigs, volumeConfig)
	}

	return volumeConfigs
}
