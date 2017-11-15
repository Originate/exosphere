package application

import (
	"path"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
)

// GenerateComposeFiles generates all docker-compose files for exosphere commands
func GenerateComposeFiles(appContext types.AppContext) error {
	buildModes := []composebuilder.BuildMode{
		composebuilder.BuildMode{
			Type:        composebuilder.BuildModeTypeLocal,
			Environment: composebuilder.BuildModeEnvironmentTest,
		},
		composebuilder.BuildMode{
			Type:        composebuilder.BuildModeTypeLocal,
			Mount:       true,
			Environment: composebuilder.BuildModeEnvironmentDevelopment,
		},
		composebuilder.BuildMode{
			Type:        composebuilder.BuildModeTypeLocal,
			Environment: composebuilder.BuildModeEnvironmentProduction,
		},
	}
	composeDir := path.Join(appContext.Location, "docker-compose")
	for _, buildMode := range buildModes {
		dockerCompose, err := composebuilder.GetApplicationDockerCompose(composebuilder.ApplicationOptions{
			AppConfig: appContext.Config,
			AppDir:    appContext.Location,
			BuildMode: buildMode,
		})
		if err != nil {
			return err
		}
		err = composebuilder.WriteYML(composeDir, buildMode.GetDockerComposeFileName(), dockerCompose)
		if err != nil {
			return err
		}
	}
	return nil
}

//GenerateTerraformFiles generates the terraform/main.tf file
func GenerateTerraformFiles(deployConfig types.DeployConfig) error {
	return terraform.GenerateFile(deployConfig)
}
