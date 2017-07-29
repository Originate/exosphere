package awsHelper

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

const secretsFile string = "secrets.tfvars"

// CreateSecretsStore creates an S3 bucket  and file object used for secrets management
func CreateSecretsStore(secretsBucket, region string) error {
	s3client := createS3client(region)

	err := createBucket(s3client, secretsBucket)
	if err != nil {
		return err
	}

	return createS3Object(s3client, nil, secretsBucket, secretsFile)
}

// ReadSecrets reads secret key value pair from remote store
func ReadSecrets(secretsBucket, region string) (types.TFString, error) {
	s3client := createS3client(region)
	err := createBucket(s3client, secretsBucket)
	if err != nil {
		return "", err
	}

	results, err := s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(secretsBucket),
		Key:    aws.String(secretsFile),
	})
	// create file if it doesn't already exist
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return "", putS3Object(s3client, nil, secretsBucket, secretsFile)
		} else {
			return "", err
		}
	}

	objectBytes, err := ioutil.ReadAll(results.Body)
	if err != nil {
		return "", err
	}

	err = results.Body.Close()
	if err != nil {
		return "", err
	}

	return types.TFString(objectBytes), err
}

// CreateSecrets creates new secret key value pair
func CreateSecrets(newSecrets map[string]string, secretsBucket, region string) error {
	s3client := createS3client(region)
	err := createBucket(s3client, secretsBucket)
	if err != nil {
		return err
	}

	err = createS3Object(s3client, nil, secretsBucket, secretsFile)
	if err != nil {
		return err
	}

	tfvars, err := ReadSecrets(secretsBucket, region)
	if err != nil {
		return err
	}

	secrets, err := ValidateAndMergeSecrets(tfvars, newSecrets)
	if err != nil {
		return err
	}

	fileBytes := bytes.NewReader([]byte(secrets))
	return putS3Object(s3client, fileBytes, secretsBucket, secretsFile)
}

// ValidateAndMergeSecrets makes sures two maps do not have conflicting keys and merges them
func ValidateAndMergeSecrets(tfvars types.TFString, newSecrets map[string]string) (types.TFString, error) {
	existingSecrets := tfvars.ToMap()
	if existingSecrets.HasConflictingKey(newSecrets) {
		return "", errors.New("new secrets have key(s) that conflict with existing secrets. Use 'exo configure update' to update existing keys")
	}

	secrets := existingSecrets.MergeSecrets(newSecrets)
	return secrets.ToTfString(), nil
}
