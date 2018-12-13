package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMKubernetesCluster_basic(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_basic(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config_raw", ""),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	location := testLocation()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(ri, location, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config_raw", ""),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantId := os.Getenv("ARM_TENANT_ID")
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(ri, location, clientId, clientSecret, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role_based_access_control.0.azure_active_directory.0.client_app_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role_based_access_control.0.azure_active_directory.0.server_app_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_admin_config_raw"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_internalNetwork(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "kubenet"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "kubenet"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.oms_agent.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.oms_agent.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "addon_profile.0.oms_agent.0.log_analytics_workspace_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.http_application_routing.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.http_application_routing.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "addon_profile.0.http_application_routing.0.http_application_routing_zone_name"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_basic(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_basic(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(rInt int, location, clientId, clientSecret string) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControl(rInt, location, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(rInt int, location, clientId, clientSecret, tenantId string) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControlAAD(rInt, location, clientId, clientSecret, tenantId)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_internalNetwork(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_internalNetwork(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworking(rInt, clientId, clientSecret, location, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingComplete(rInt, clientId, clientSecret, location, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworking(rInt, clientId, clientSecret, location, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingComplete(rInt, clientId, clientSecret, location, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileOMS(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileRouting(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}
