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
func PushApplicationImages(options PushApplicationImagesOptions) (map[string]string, error) {
	config := aws.CreateAwsConfig(options.DeployConfig.AwsConfig)
	session := session.Must(session.NewSession())
	ecrClient := ecr.New(session, config)
	ecrAuth, err := aws.GetECRCredentials(ecrClient)
	if err != nil {
		return nil, err
	}
	dockerCompose, err := tools.GetDockerCompose(path.Join(options.DockerComposeDir, "docker-compose.yml"))
	if err != nil {
		return nil, err
	}
	imagesMap, err := GetImageNames(options.DeployConfig, options.DockerComposeDir, dockerCompose)
	if err != nil {
		return nil, err
	}
	serviceData := options.DeployConfig.AppConfig.GetServiceData()
	for serviceRole, imageName := range imagesMap {
		taggedImage, err := PushServiceImage(PushServiceImageOptions{
			DeployConfig:     options.DeployConfig,
			DockerComposeDir: options.DockerComposeDir,
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
