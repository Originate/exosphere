package aws

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// CreateRepository creates a repository with the given name if one does not exist
// returns the repositoryURI
func CreateRepository(ecrClient *ecr.ECR, repositoryName string) (string, error) {
	repositoryURI, err := getRepositoryURI(ecrClient, repositoryName)
	if err != nil {
		return "", err
	}
	if repositoryURI != "" {
		return repositoryURI, nil
	}
	result, err := ecrClient.CreateRepository(&ecr.CreateRepositoryInput{
		RepositoryName: aws.String(repositoryName),
	})
	if err != nil {
		return "", err
	}
	return *result.Repository.RepositoryUri, nil
}

// GetECRCredentials returns base64 encoded ECR auth object
func GetECRCredentials(ecrClient *ecr.ECR) (string, error) {
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

// Helpers

// retrieves encoded ECR credentails (in the format username:password)
// and returns them as separate strings
func getEcrAuth(ecrClient *ecr.ECR) (string, string, error) {
	result, err := ecrClient.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", "", err
	}
	str := *result.AuthorizationData[0].AuthorizationToken
	decodedAuth, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", "", err
	}
	decodedAuthArgs := strings.Split(string(decodedAuth), ":")
	return decodedAuthArgs[0], decodedAuthArgs[1], nil
}

func getRepositoryURI(ecrClient *ecr.ECR, repositoryName string) (string, error) {
	result, err := ecrClient.DescribeRepositories(&ecr.DescribeRepositoriesInput{})
	if err != nil {
		return "", err
	}
	for _, repositoryInfo := range result.Repositories {
		if *repositoryInfo.RepositoryName == repositoryName {
			return *repositoryInfo.RepositoryUri, nil
		}
	}
	return "", nil
}
