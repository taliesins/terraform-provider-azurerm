package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"azurerm": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	variables := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_SUBSCRIPTION_ID",
		"ARM_TENANT_ID",
		"ARM_TEST_LOCATION",
		"ARM_TEST_LOCATION_ALT",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func testLocation() string {
	return os.Getenv("ARM_TEST_LOCATION")
}

func testAltLocation() string {
	return os.Getenv("ARM_TEST_LOCATION_ALT")
}

func testArmEnvironmentName() string {
	envName, exists := os.LookupEnv("ARM_ENVIRONMENT")
	if !exists {
		envName = "public"
	}

	return envName
}

func testArmEnvironment() (*azure.Environment, error) {
	envName := testArmEnvironmentName()
	return authentication.DetermineEnvironment(envName)
}

func testGetAzureConfig(t *testing.T) *authentication.Config {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Integration test skipped unless env '%s' set", resource.TestEnvVar))
		return nil
	}

	environment := testArmEnvironmentName()

	builder := authentication.Builder{
		SubscriptionID: os.Getenv("ARM_SUBSCRIPTION_ID"),
		ClientID:       os.Getenv("ARM_CLIENT_ID"),
		TenantID:       os.Getenv("ARM_TENANT_ID"),
		ClientSecret:   os.Getenv("ARM_CLIENT_SECRET"),
		Environment:    environment,

		// we intentionally only support Client Secret auth for tests (since those variables are used all over)
		SupportsClientSecretAuth: true,
	}
	config, err := builder.Build()
	if err != nil {
		t.Fatalf("Error building ARM Client: %+v", err)
		return nil
	}

	return config
}
