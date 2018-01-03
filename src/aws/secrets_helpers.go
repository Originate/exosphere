package aws

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

const secretsFile string = "secrets.json"

// ReadSecrets reads secret key value pair from remote store
// It creates an empty secrets store if necessary
func ReadSecrets(awsConfig types.AwsConfig) (types.Secrets, error) {
	s3client := createS3client(awsConfig)
	err := createS3Object(s3client, strings.NewReader("{}"), awsConfig.BucketName, secretsFile)
	if err != nil {
		return nil, err
	}
	results, err := s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(awsConfig.BucketName),
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
	util.Merge(existingSecrets, newSecrets)
	return writeSecrets(existingSecrets, awsConfig)
}

// DeleteSecrets deletes a list of secrets provided their keys. Ignores them if they don't exist
func DeleteSecrets(existingSecrets types.Secrets, secretKeys []string, awsConfig types.AwsConfig) error {
	newSecrets := existingSecrets.Delete(secretKeys)
	return writeSecrets(newSecrets, awsConfig)
}

func writeSecrets(secrets types.Secrets, awsConfig types.AwsConfig) error {
	s3client := createS3client(awsConfig)
	secretsString, err := json.Marshal(secrets)
	if err != nil {
		return errors.Wrap(err, "cannot marshal secrets map into JSON string")
	}
	fileBytes := bytes.NewReader(secretsString)
	return putS3Object(s3client, fileBytes, awsConfig.BucketName, secretsFile)
}
