package deployer

import (
	"path"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// PushApplicationImages pushes all the docker images for the application to ECR
func PushApplicationImages(deployConfig deploy.Config) (map[string]string, error) {
	config := aws.CreateAwsConfig(deployConfig.GetAwsOptions())
	session := session.Must(session.NewSession())
	ecrClient := ecr.New(session, config)
	ecrAuth, err := aws.GetECRCredentials(ecrClient)
	if err != nil {
		return nil, err
	}
	dockerCompose, err := tools.GetDockerCompose(path.Join(deployConfig.AppContext.GetDockerComposeDir(), types.LocalProductionComposeFileName))
	if err != nil {
		return nil, err
	}
	imagesMap := GetServiceImageNames(deployConfig, dockerCompose)
	serviceData := deployConfig.AppContext.Config.Services
	for serviceRole, imageName := range imagesMap {
		taggedImage, err := PushImage(PushImageOptions{
			DeployConfig:    deployConfig,
			EcrAuth:         ecrAuth,
			EcrClient:       ecrClient,
			ImageName:       imageName,
			ServiceRole:     serviceRole,
			ServiceLocation: serviceData[serviceRole].Location,
		})
		if err != nil {
			return nil, err
		}
		imagesMap[serviceRole] = taggedImage
	}
	return imagesMap, nil
}
