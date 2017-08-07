package aws

import (
	"encoding/base64"
	"encoding/json"
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
// Create the image repository on ECR if it doesn't exist
// Tags application images
// Retrieves ECR credentials
// Pushes images to ECR
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
		taggedImage := fmt.Sprintf("%s:%s", repositoryURI, appConfig.Version)
		err = docker.TagImage(dockerClient, imageName, taggedImage)
		if err != nil {
			return err
		}
		encodedAuth, err := getECRCredentials(ecrClient)
		if err != nil {
			return err
		}
		err = docker.PushImage(dockerClient, taggedImage, encodedAuth)
		if err != nil {
			return err
		}
	}
	return nil
}

func getECRCredentials(ecrClient *ecr.ECR) (string, error) {
	registryUser, registryPass, err := getEcrAuth(ecrClient)
	if err != nil {
		return "", err
	}
	authObj := map[string]string{
		"username": registryUser,
		"password": registryPass,
	}
	json, err := json.Marshal(authObj)
	if err != nil {
		return "", err
	}
	encodedAuth := base64.StdEncoding.EncodeToString(json)
	return encodedAuth, nil
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
