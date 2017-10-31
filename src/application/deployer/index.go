package deployer

import (
	"fmt"
	"path"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	prompt "github.com/kofalt/go-prompt"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	err := validateConfigs(deployConfig)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Setting up AWS account...")
	err = aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Pushing Docker images to ECR...")
	dockerConfigs, err := composebuilder.GetApplicationDockerConfigs(composebuilder.ApplicationOptions{
		AppConfig: deployConfig.AppContext.Config,
		AppDir:    deployConfig.AppContext.Location,
		BuildMode: composebuilder.BuildMode{
			Type: composebuilder.BuildModeTypeDeploy,
		},
		HomeDir: deployConfig.HomeDir,
	})
	if err != nil {
		return err
	}
	dockerComposeDir := path.Join(deployConfig.AppContext.Location, "tmp")
	err = composebuilder.WriteYML(dockerComposeDir, dockerConfigs)
	if err != nil {
		return err
	}
	imagesMap, err := PushApplicationImages(deployConfig, dockerComposeDir)
	if err != nil {
		return err
	}

	prevTerraformFileContents, err := terraform.ReadTerraformFile(deployConfig)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Generating Terraform files...")
	err = terraform.GenerateFile(deployConfig, imagesMap)
	if err != nil {
		return err
	}

	if deployConfig.DeployServicesOnly {
		err = terraform.CheckTerraformFile(deployConfig, prevTerraformFileContents)
		if err != nil {
			return err
		}
	}

	return deployApplication(deployConfig, imagesMap, prevTerraformFileContents)
}

func validateConfigs(deployConfig types.DeployConfig) error {
	fmt.Fprintln(deployConfig.Writer, "Validating application configuration...")
	err := deployConfig.AppContext.Config.Production.ValidateFields()
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Validating service configurations...")
	protectionLevels := deployConfig.AppContext.Config.GetServiceProtectionLevels()
	serviceData := deployConfig.AppContext.Config.GetServiceData()
	for serviceRole, serviceConfig := range deployConfig.ServiceConfigs {
		err = serviceConfig.Production.ValidateFields(serviceData[serviceRole].Location, protectionLevels[serviceRole])
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(deployConfig.Writer, "Validating application dependencies...")
	validatedDependencies := map[string]string{}
	for _, dependency := range deployConfig.AppContext.Config.Production.Dependencies {
		err = dependency.ValidateFields()
		if err != nil {
			return err
		}
		validatedDependencies[dependency.Name] = dependency.Version
	}

	fmt.Fprintln(deployConfig.Writer, "Validating service dependencies...")
	for _, serviceConfig := range deployConfig.ServiceConfigs {
		for _, dependency := range serviceConfig.Production.Dependencies {
			if validatedDependencies[dependency.Name] == "" {
				err = dependency.ValidateFields()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func deployApplication(deployConfig types.DeployConfig, imagesMap map[string]string, prevTerraformFileContents []byte) error {
	fmt.Fprintln(deployConfig.Writer, "Retrieving remote state...")
	err := terraform.RunInit(deployConfig)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Retrieving secrets...")
	secrets, err := aws.ReadSecrets(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	applyPlan := true
	if !deployConfig.DeployServicesOnly {
		fmt.Fprintln(deployConfig.Writer, "Planning deployment...")
		err = terraform.RunPlan(deployConfig, secrets, imagesMap)
		if err != nil {
			return err
		}
		applyPlan = prompt.Confirm("Do you want to apply this plan? (y/n)")
	}

	if applyPlan {
		fmt.Fprintln(deployConfig.Writer, "Applying changes...")
		return terraform.RunApply(deployConfig, secrets, imagesMap)
	}
	fmt.Fprintln(deployConfig.Writer, "Abandoning deployment...reverting 'terraform/main.tf' file.")
	return terraform.WriteTerraformFile(string(prevTerraformFileContents), deployConfig.TerraformDir)
}
