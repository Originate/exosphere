package deployer

import (
	"path"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/Originate/exosphere/src/types"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// PushApplicationImagesOptions is the options to PushApplicationImages
type PushApplicationImagesOptions struct {
	DeployConfig     types.DeployConfig
	DockerComposeDir string
}

// PushApplicationImages pushes all the docker images for the application to ECR
func PushApplicationImages(deployConfig types.DeployConfig, dockerComposeDir string) (map[string]string, error) {
	config := aws.CreateAwsConfig(deployConfig.AwsConfig)
	session := session.Must(session.NewSession())
	ecrClient := ecr.New(session, config)
	ecrAuth, err := aws.GetECRCredentials(ecrClient)
	if err != nil {
		return nil, err
	}
	dockerCompose, err := tools.GetDockerCompose(path.Join(dockerComposeDir, "docker-compose.yml"))
	if err != nil {
		return nil, err
	}
	imagesMap, err := GetImageNames(deployConfig, dockerComposeDir, dockerCompose)
	if err != nil {
		return nil, err
	}
	serviceData := deployConfig.AppConfig.GetServiceData()
	for serviceRole, imageName := range imagesMap {
		taggedImage, err := PushImage(PushImageOptions{
			DeployConfig:     deployConfig,
			DockerComposeDir: dockerComposeDir,
			EcrAuth:          ecrAuth,
			EcrClient:        ecrClient,
			ImageName:        imageName,
			ServiceLocation:  serviceData[serviceRole].Location,
		})
		if err != nil {
			return nil, err
		}
		imagesMap[serviceRole] = taggedImage
	}
	return imagesMap, nil
}
