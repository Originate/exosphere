package aws

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

const secretsFile string = "secrets.json"

// CreateSecretsStore creates an S3 bucket  and file object used for secrets management
func CreateSecretsStore(awsConfig types.AwsConfig) error {
	s3client := createS3client(awsConfig.Region)
	return createS3Object(s3client, strings.NewReader("{}"), awsConfig.SecretsBucket, secretsFile)
}

// ReadSecrets reads secret key value pair from remote store
func ReadSecrets(awsConfig types.AwsConfig) (types.Secrets, error) {
	s3client := createS3client(awsConfig.Region)
	err := createS3Object(s3client, strings.NewReader("{}"), awsConfig.SecretsBucket, secretsFile)
	if err != nil {
		return nil, err
	}
	results, err := s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(awsConfig.SecretsBucket),
		Key:    aws.String(secretsFile),
	})
	if err != nil {
		return nil, err
	}
	objectBytes, err := ioutil.ReadAll(results.Body)
	if err != nil {
		return nil, err
	}
	err = results.Body.Close()
	if err != nil {
		return nil, err
	}
	secrets := types.Secrets{}
	err = json.Unmarshal(objectBytes, &secrets)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal secrets into map")
	}
	return secrets, nil
}

// MergeAndWriteSecrets merges two secret maps and writes them to s3
// Overwrites existingSecrets's values if the are conflicting keys
func MergeAndWriteSecrets(existingSecrets, newSecrets types.Secrets, awsConfig types.AwsConfig) error {
	secrets := existingSecrets.Merge(newSecrets)
	return writeSecrets(secrets, awsConfig)
}

// DeleteSecrets deletes a list of secrets provided their keys. Ignores them if they don't exist
func DeleteSecrets(secretKeys []string, awsConfig types.AwsConfig) error {
	secrets, err := ReadSecrets(awsConfig)
	if err != nil {
		return err
	}
	newSecrets := secrets.Delete(secretKeys)
	return writeSecrets(newSecrets, awsConfig)
}

func writeSecrets(secrets types.Secrets, awsConfig types.AwsConfig) error {
	s3client := createS3client(awsConfig.Region)
	secretsString, err := json.Marshal(secrets)
	if err != nil {
		return errors.Wrap(err, "cannot marshal secrets map into JSON string")
	}
	fileBytes := bytes.NewReader([]byte(secretsString))
	return putS3Object(s3client, fileBytes, awsConfig.SecretsBucket, secretsFile)
}
