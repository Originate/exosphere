package deployer

import (
	"path"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/Originate/exosphere/src/types"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// PushApplicationImages pushes all the docker images for the application to ECR
func PushApplicationImages(deployConfig types.DeployConfig, dockerComposeDir string, buildMode composebuilder.BuildMode) (map[string]string, error) {
	config := aws.CreateAwsConfig(deployConfig.AwsConfig)
	session := session.Must(session.NewSession())
	ecrClient := ecr.New(session, config)
	ecrAuth, err := aws.GetECRCredentials(ecrClient)
	if err != nil {
		return nil, err
	}
	dockerCompose, err := tools.GetDockerCompose(path.Join(dockerComposeDir, buildMode.GetDockerComposeFileName()))
	if err != nil {
		return nil, err
	}
	imagesMap, err := GetImageNames(deployConfig, dockerComposeDir, dockerCompose)
	if err != nil {
		return nil, err
	}
	serviceData := deployConfig.AppContext.Config.Services
	for serviceRole, imageName := range imagesMap {
		taggedImage, err := PushImage(PushImageOptions{
			DeployConfig:     deployConfig,
			DockerComposeDir: dockerComposeDir,
			EcrAuth:          ecrAuth,
			EcrClient:        ecrClient,
			ImageName:        imageName,
			ServiceRole:      serviceRole,
			ServiceLocation:  serviceData[serviceRole].Location,
			BuildMode:        buildMode,
		})
		if err != nil {
			return nil, err
		}
		imagesMap[serviceRole] = taggedImage
	}
	return imagesMap, nil
}
