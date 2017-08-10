package aws

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker"
	"github.com/Originate/exosphere/src/types"
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
// Returns a map from service name to image name on ECR
func PushImages(deployConfig types.DeployConfig, dockerComposePath string) (map[string]string, error) {
	config := aws.NewConfig().WithRegion(deployConfig.AwsConfig.Region)
	session := session.Must(session.NewSession())
	ecrClient := ecr.New(session, config)
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	dockerCompose, err := docker.GetDockerCompose(dockerComposePath)
	if err != nil {
		return nil, err
	}
	imagesMap, err := getImageNames(deployConfig, filepath.Dir(dockerComposePath), dockerCompose)
	if err != nil {
		return nil, err
	}
	for serviceName, imageName := range imagesMap {
		deployConfig.LogChannel <- fmt.Sprintf("Pushing image for: %s...\n\n", serviceName)
		repositoryName, version := getRepositoryConfig(deployConfig, imageName)
		repositoryURI, err := createRepository(ecrClient, repositoryName)
		if err != nil {
			return nil, err
		}
		taggedImage := fmt.Sprintf("%s:%s", repositoryURI, version)
		err = docker.TagImage(dockerClient, imageName, taggedImage)
		if err != nil {
			return nil, err
		}
		imagesMap[serviceName] = taggedImage
		encodedAuth, err := getECRCredentials(ecrClient)
		if err != nil {
			return nil, err
		}
		err = docker.PushImage(dockerClient, taggedImage, encodedAuth)
		if err != nil {
			return nil, err
		}
	}
	return imagesMap, nil
}

// returns base64 encoded ECR auth object
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

// returns a mapping from service/dependency names to image name on the user's machine
func getImageNames(deployConfig types.DeployConfig, dockerComposeDir string, dockerCompose types.DockerCompose) (map[string]string, error) {
	images := map[string]string{}
	// get service image names
	for _, serviceName := range deployConfig.AppConfig.GetServiceNames() {
		dockerConfig := dockerCompose.Services[serviceName]
		images[serviceName] = buildImageName(dockerConfig, dockerComposeDir, serviceName)
	}
	// get dependency image names
	serviceConfigs, err := config.GetServiceConfigs(deployConfig.AppDir, deployConfig.AppConfig)
	if err != nil {
		return nil, err
	}
	dependencies := config.GetAllBuiltDependencies(deployConfig.AppConfig, serviceConfigs, deployConfig.AppDir, deployConfig.HomeDir)
	for dependencyName, dependency := range dependencies {
		dockerConfig, err := dependency.GetDockerConfig()
		if err != nil {
			return nil, err
		}
		images[dependencyName] = buildImageName(dockerConfig, dockerComposeDir, dependencyName)
	}
	return images, nil
}

// returns image name as it appears on the user's machine
func buildImageName(dockerConfig types.DockerConfig, dockerComposeDir, serviceName string) string {
	if dockerConfig.Image != "" {
		return dockerConfig.Image
	}
	return fmt.Sprintf("%s_%s", filepath.Base(dockerComposeDir), serviceName)
}

// returns an image with version tag if applicable. uses the application version otherwise
func getRepositoryConfig(deployConfig types.DeployConfig, imageName string) (string, string) {
	config := strings.Split(imageName, ":")
	repositoryName := config[0]
	var version string
	if len(config) > 1 {
		version = config[1]
	} else {
		version = deployConfig.AppConfig.Version
	}
	return repositoryName, version
}
