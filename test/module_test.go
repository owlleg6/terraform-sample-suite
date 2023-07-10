package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/operationalinsights/mgmt/operationalinsights"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var (
	subscriptionId         = os.Getenv("ARM_SUBSCRIPTION_ID")
	uniqueEnv              = random.UniqueId()
	uniqueSuffix           = random.UniqueId()
	terraformDefaultConfig = &terraform.Options{
		TerraformDir: "../terraform",
	}
)

// Creates Resource Group with default values
func TestCreateResouceGroup(t *testing.T) {
	// Copies default Terraform config
	terraformOptions := terraformDefaultConfig

	// Runs Terraform Init and Apply
	terraform.InitAndApply(t, terraformOptions)
}

// Recreates Resource Group with suffix in name
func TestRecreateWithSuffix(t *testing.T) {
	// Copies default Terraform config
	terraformOptions := terraformDefaultConfig

	// Sets Terraform variable 'suffix' to new value
	terraformOptions.Vars = map[string]interface{}{
		"suffix": uniqueSuffix,
		"env":    uniqueEnv,
	}

	// Runs Terraform Apply only
	terraform.Apply(t, terraformOptions)

	// Stores value of 'name' parameter of from Terraform outputs
	resourceGroupName := terraform.OutputMap(t, terraformOptions, "resource_group")["name"]

	// Checks if created resource exists in Azure under correct name (from outputs)
	exists := azure.ResourceGroupExists(t, resourceGroupName, subscriptionId)

	// If it can't find such resource test will throw an Error
	assert.True(t,
		exists,
		"Resource group with suffix does not exist",
	)
}

// Creates Log Analytics Workspace
func TestCreateLogAnalyticsWorkspace(t *testing.T) {
	// Copies default Terraform config
	terraformOptions := terraformDefaultConfig

	// Sets Terraform variable 'suffix' to new value
	terraformOptions.Vars = map[string]interface{}{
		"suffix":                   uniqueSuffix,
		"env":                      uniqueEnv,
		"log_analytics_ws_enabled": true,
	}

	// Runs Terraform Apply only
	terraform.InitAndApply(t, terraformOptions)

	// Stores value of 'name' parameter of from Terraform outputs
	resourceGroupName := terraform.OutputMap(t, terraformOptions, "resource_group")["name"]

	// Stores value of 'name' parameter of from Terraform outputs
	logAnalyticsName := terraform.OutputMap(t, terraformOptions, "log_analytics_ws")["name"]

	// Checks if created resource exists in Azure with correct name (from outputs)
	exists := azure.LogAnalyticsWorkspaceExists(
		t, logAnalyticsName, resourceGroupName, subscriptionId,
	)

	// If it cant find such resource test will throw an Error
	assert.True(t,
		exists,
		"Log Analytics Workspace does not exists",
	)
}

// Recreates resouces with fully custom name
func TestRecreateWithCustomNames(t *testing.T) {
	// Pointer example. All changes would be also saved in original variable
	terraformOptions := *terraformDefaultConfig

	// Provides values to variables that define custom names for resources
	terraformOptions.Vars = map[string]interface{}{
		"custom_resource_group_name":  fmt.Sprintf("resource-group-%v-%v", uniqueEnv, uniqueSuffix),
		"log_analytics_ws_enabled":    true,
		"custom_workspace_name":       fmt.Sprintf("log-analytics-ws-%v-%v", uniqueEnv, uniqueSuffix),
		"analytics_retention_in_days": 15,
	}

	// In this case, test will run Apply and then Plan,
	// if there will be changes in new Plan this test will throw an Error
	//terraform.ApplyAndIdempotent(t, terraformDefaultConfig)
	terraform.InitAndApply(t, terraformDefaultConfig)

	// Stores value of 'name' parameter of from Terraform outputs
	resourceGroupName := terraform.OutputMap(t, terraformDefaultConfig, "resource_group")["name"]

	// Stores value of 'name' parameter of from Terraform outputs
	logAnalyticsName := terraform.OutputMap(t, terraformDefaultConfig, "log_analytics_ws")["name"]

	// Checks if created resource exists in Azure with correct name (from outputs)
	existsRg := azure.ResourceGroupExists(
		t, resourceGroupName, subscriptionId,
	)

	// Checks if created resource exists in Azure with correct name (from outputs)
	existsWs := azure.LogAnalyticsWorkspaceExists(
		t, logAnalyticsName, resourceGroupName, subscriptionId,
	)

	// If it cant find such resource under specified name -  test will throw an Error
	assert.True(t,
		existsRg,
		"Resource group with custom name does not exists",
	)

	assert.True(t,
		existsWs,
		"Log Analytics Workspace with custom name does not exists",
	)
}

// Validates configuration of Log Analytics Workspaces
func TestLogAnalyticsWsConfigTesting(t *testing.T) {
	// Stores value of 'name' parameter of from Terraform outputs
	resourceGroupName := terraform.OutputMap(t, terraformDefaultConfig, "resource_group")["name"]

	// Stores value of 'name' parameter of from Terraform outputs
	logAnalyticsName := terraform.OutputMap(t, terraformDefaultConfig, "log_analytics_ws")["name"]

	// Gets all infromation about created Log Analytics Workspaces
	ws := azure.GetLogAnalyticsWorkspace(t, logAnalyticsName, resourceGroupName, subscriptionId)

	// Checks if Workspace successfully created
	assert.EqualValues(t,
		operationalinsights.WorkspaceEntityStatusSucceeded,
		ws.WorkspaceProperties.ProvisioningState,
		"Workspace creation failed",
	)

	// Checks if public network access to workspace enabled
	assert.EqualValues(t,
		operationalinsights.Enabled,
		ws.WorkspaceProperties.PublicNetworkAccessForIngestion,
		"Public Access is disabled on Workspace",
	)

	// Checks if Workspace Log retention duration is greater or equal to provided by Terraform
	assert.EqualValues(t,
		int32(15), //terraformDefaultConfig.Vars["analytics_retention_in_days"],
		ws.WorkspaceProperties.RetentionInDays,
		"Workspace log retention duration is not valid",
	)
}

// Destroys all Terraform resources
func TestDestroyResources(t *testing.T) {
	terraform.Destroy(t, terraformDefaultConfig)
}

// Helper function. Removes tfstate files from Terraform directory
func TestHelperRemovesTerraformFiles(t *testing.T) {
	dirPath := terraformDefaultConfig.TerraformDir

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through each file
	for _, f := range files {
		// If the file ends with ".tfstate", delete it
		if strings.HasSuffix(f.Name(), ".tfstate.backup") {
			err = os.Remove(filepath.Join(dirPath, f.Name()))
			if err != nil {
				log.Printf("Failed to delete file: %s", f.Name())
			}
		}
	}
}
