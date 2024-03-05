package test

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"path/filepath"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformOSDUInstaller(t *testing.T) {
	
	// List of environment variables to check
	envVariables := []string{"TF_MODULE_DIR", "TF_MODULE_NAME", "TERRAFORM_ENVIRONMENT", "DIST_DIR"}
	// Loop through the list of environment variables
	for _, envVar := range envVariables {
		value, exists := os.LookupEnv(envVar)
		if exists {
			fmt.Printf("Environment variable %s is set with value: %s\n", envVar, value)
		} else {
			fmt.Printf("Environment variable %s is not set\n", envVar)
		}
	}

	tfModuleDir := os.Getenv("TF_MODULE_DIR")
	tfModuleName := os.Getenv("TF_MODULE_NAME")
	tfEnvironment := os.Getenv("TERRAFORM_ENVIRONMENT")
	distDir := os.Getenv("DIST_DIR")
	tfInitUpgrade := true

	
	tfVarFile := fmt.Sprintf("osdu/%s/%s/%s.tfvars", distDir, tfEnvironment, tfModuleName)
	tfConfFile := fmt.Sprintf("osdu/%s/%s/%s.conf", distDir, tfEnvironment, tfModuleName)
	tfDir := fmt.Sprintf("osdu/deployment/%s", tfModuleDir)

	//terratest requires absolute patf of tfvars file
	absolutePathTfVar, err := filepath.Abs(tfVarFile)
	if err != nil {
		fmt.Println("Error:", err)
	}

	tfVarFileAbsPath := []string{absolutePathTfVar}
	
	//Convert conf file into map of strings
	file, err := os.Open(tfConfFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	tfConf := make(map[string]interface{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			tfConf[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning the file:", err)
		return
	}

	//Run Terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: tfDir,
		BackendConfig: tfConf,
		Upgrade: tfInitUpgrade, 
		VarFiles: tfVarFileAbsPath,
	})

	// // Clean up resources with "terraform destroy" at the end of the test.
	// defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

}
