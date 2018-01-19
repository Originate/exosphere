package terraform

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/util"
)

func getSharedVarMap(deployConfig deploy.Config, secrets types.Secrets) map[string]string {
	varMap := map[string]string{}
	util.Merge(varMap, secrets)
	util.Merge(varMap, getAwsVarMap(deployConfig))
	varMap["application_url"] = deployConfig.AppContext.Config.Remote.Environments[deployConfig.RemoteEnvironmentID].URL
	varMap["env"] = deployConfig.RemoteEnvironmentID
	return varMap
}

func getAwsVarMap(deployConfig deploy.Config) map[string]string {
	remoteEnvironment := deployConfig.GetRemoteEnvironment()
	return map[string]string{
		"aws_profile":             deployConfig.AwsProfile,
		"aws_region":              remoteEnvironment.Region,
		"aws_account_id":          remoteEnvironment.AccountID,
		"aws_ssl_certificate_arn": remoteEnvironment.SslCertificateArn,
	}
}

// convert an env var key pair in the format of a task definition
// marshals a map[string]string object twice so that it can be escaped properly
// and passed as a command line flag, then properly decoded in Terraform
func createEnvVarString(envVars map[string]string) (string, error) {
	terraformEnvVars := []map[string]string{}
	for k, v := range envVars {
		envVarPair := map[string]string{
			"name":  k,
			"value": v,
		}
		terraformEnvVars = append(terraformEnvVars, envVarPair)
	}
	envVarsJSON, err := json.Marshal(terraformEnvVars)
	if err != nil {
		return "", err
	}
	return string(envVarsJSON), nil
}

func writeVarFile(varMap map[string]string, terraformDir, remoteEnvironmentID string) error {
	jsonVarMap, err := json.MarshalIndent(varMap, "", "  ")
	if err != nil {
		return err
	}
	return WriteToNestedTerraformDir(string(jsonVarMap), fmt.Sprintf("%s.tfvars", remoteEnvironmentID), terraformDir)
}
