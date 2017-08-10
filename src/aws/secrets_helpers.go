package aws

import (
	"bytes"
	"io/ioutil"

	"github.com/Originate/exosphere/src/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const secretsFile string = "secrets.tfvars"

// CreateSecretsStore creates an S3 bucket  and file object used for secrets management
func CreateSecretsStore(awsConfig types.AwsConfig) error {
	s3client := createS3client(awsConfig.Region)
	return createS3Object(s3client, nil, awsConfig.SecretsBucket, secretsFile)
}

// ReadSecrets reads secret key value pair from remote store
func ReadSecrets(awsConfig types.AwsConfig) (string, error) {
	s3client := createS3client(awsConfig.Region)
	err := createS3Object(s3client, nil, awsConfig.SecretsBucket, secretsFile)
	if err != nil {
		return "", err
	}
	results, err := s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(awsConfig.SecretsBucket),
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

// MergeAndWriteSecrets merges two secret maps and writes them to s3
// Overwrites existingSecrets's values if the are conflicting keys
func MergeAndWriteSecrets(existingSecrets, newSecrets types.Secrets, awsConfig types.AwsConfig) error {
	secrets := existingSecrets.Merge(newSecrets)
	return writeSecrets(secrets, awsConfig)
}

// DeleteSecrets deletes a list of secrets provided their keys. Ignores them if they don't exist
func DeleteSecrets(secretKeys []string, awsConfig types.AwsConfig) error {
	tfvars, err := ReadSecrets(awsConfig)
	if err != nil {
		return err
	}
	secrets := types.NewSecrets(tfvars).Delete(secretKeys)
	return writeSecrets(secrets, awsConfig)
}

func writeSecrets(secrets types.Secrets, awsConfig types.AwsConfig) error {
	s3client := createS3client(awsConfig.Region)
	fileBytes := bytes.NewReader([]byte(secrets.TfString()))
	return putS3Object(s3client, fileBytes, awsConfig.SecretsBucket, secretsFile)
}
