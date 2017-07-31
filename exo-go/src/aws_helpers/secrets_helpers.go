package awsHelper

import (
	"bytes"
	"io/ioutil"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const secretsFile string = "secrets.tfvars"

// CreateSecretsStore creates an S3 bucket  and file object used for secrets management
func CreateSecretsStore(secretsBucket, region string) error {
	s3client := createS3client(region)
	return createS3Object(s3client, nil, secretsBucket, secretsFile)
}

// ReadSecrets reads secret key value pair from remote store
func ReadSecrets(secretsBucket, region string) (string, error) {
	s3client := createS3client(region)
	err := createS3Object(s3client, nil, secretsBucket, secretsFile)
	if err != nil {
		return "", err
	}

	results, err := s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(secretsBucket),
		Key:    aws.String(secretsFile),
	})
	if err != nil {
		return "", err
	}

	objectBytes, err := ioutil.ReadAll(results.Body)
	if err != nil {
		return "", err
	}

	err = results.Body.Close()
	if err != nil {
		return "", err
	}

	return string(objectBytes), err
}

// CreateSecrets creates new secret key value pair
func CreateSecrets(newSecrets map[string]string, secretsBucket, region string) error {
	tfvars, err := ReadSecrets(secretsBucket, region)
	if err != nil {
		return err
	}

	secrets, err := types.NewSecrets(tfvars).ValidateAndMerge(newSecrets)
	if err != nil {
		return err
	}

	s3client := createS3client(region)
	fileBytes := bytes.NewReader([]byte(secrets.TfString()))
	return putS3Object(s3client, fileBytes, secretsBucket, secretsFile)
}
