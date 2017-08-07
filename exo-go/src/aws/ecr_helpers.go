package aws

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/docker"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/moby/moby/client"
)

// PushImages does the following:
// Logs into ECR registrar
// Create the image repository on ECR if it doesn't exist
// Tags image
// Pushes image to ECR
func PushImages(appConfig types.AppConfig, dockerComposePath, region string) error {
	config := aws.NewConfig().WithRegion(region)
	session := session.Must(session.NewSession())
	ecrClient := ecr.New(session, config)
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	dockerCompose, err := docker.GetDockerCompose(dockerComposePath)
	if err != nil {
		return err
	}
	for _, imageName := range getImageNames(filepath.Dir(dockerComposePath), dockerCompose) {
		fmt.Printf("Pushing image: %s...\n\n", imageName)
		repositoryName := strings.Split(imageName, ":")[0]
		repositoryURI, err := createRepository(ecrClient, repositoryName)
		if err != nil {
			return err
		}
		taggedImage, err := tagImage(dockerClient, imageName, repositoryURI, appConfig.Version)
		if err != nil {
			return err
		}
		registryUser, registryPass, err := getEcrAuth(ecrClient)
		if err != nil {
			return err
		}
		err = docker.PushImage(dockerClient, taggedImage, registryUser, registryPass)
		if err != nil {
			return err
		}
	}
	return nil
}

func tagImage(dockerClient *client.Client, imageName, repositoryURI, version string) (string, error) {
	taggedImage := fmt.Sprintf("%s:%s", repositoryURI, version)
	return taggedImage, docker.TagImage(dockerClient, imageName, taggedImage)
}

func getImageNames(dockerComposeDir string, dockerCompose types.DockerCompose) []string {
	images := []string{}
	for serviceName, dockerConfig := range dockerCompose.Services {
		if dockerConfig.Image != "" {
			images = append(images, dockerConfig.Image)
		} else {
			imageName := fmt.Sprintf("%s_%s", filepath.Base(dockerComposeDir), serviceName)
			images = append(images, imageName)
		}
	}
	return images
}
