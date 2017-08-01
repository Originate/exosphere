package aws

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

// MergeAndPutSecrets merges two secret maps and writes them to s3
// Overwrites existingSecrets's values if the are conflicting keys
func MergeAndWriteSecrets(existingSecrets, newSecrets types.Secrets, secretsBucket, region string) error {
	secrets := existingSecrets.Merge(newSecrets)
	return writeSecrets(secrets, secretsBucket, region)
}

// DeleteSecrets deletes a list of secrets provided their key
func DeleteSecrets(secretKeys []string, secretsBucket, region string) error {
	tfvars, err := ReadSecrets(secretsBucket, region)
	if err != nil {
		return err
	}
	secrets := types.NewSecrets(tfvars).DeleteSecrets(secretKeys)
	return writeSecrets(secrets, secretsBucket, region)
}

func writeSecrets(secrets types.Secrets, secretsBucket, region string) error {
	s3client := createS3client(region)
	fileBytes := bytes.NewReader([]byte(secrets.TfString()))
	return putS3Object(s3client, fileBytes, secretsBucket, secretsFile)
}
