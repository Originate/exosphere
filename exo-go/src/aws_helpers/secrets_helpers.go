package awsHelper

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// CreateSecretsStore creates an S3 bucket used for secrets management
func CreateSecretsStore(bucketName, region string) error {
	s3client := createS3client(region)

	hasBucket, err := hasBucket(s3client, bucketName)
	if err != nil {
		return err
	}

	if !hasBucket {
		_, err = s3client.CreateBucket(&s3.CreateBucketInput{Bucket: &bucketName})
		if err != nil {
			return err
		}
	}

	return err
}

//TODO
func ReadSecrets() error {

}

// CreateSecrets creates new secret key value pair
func CreateSecrets(secrets map[string]string, secretsBucket, region string) error {
	//TODO: first read secrets file
	return putSecretsFile(secrets, secretsBucket, region)
}

// takes a map from {"secret_name": "secret_value"} and
// creates a .tfvars file and uploads it to s3
func putSecretsFile(secrets map[string]string, secretsBucket, region string) error {
	tfvars := ""
	for key, value := range secrets {
		tfvars += fmt.Sprintf("%s=\"%s\"\n", key, value)
	}

	s3client := createS3client(region)
	hasBucket, err := hasBucket(s3client, secretsBucket)
	if err != nil {
		return err
	}
	if !hasBucket {
		return errors.New("Secrets bucket not found. Run 'exo configure' to initialize bucket.")
	}

	fileBytes := bytes.NewReader([]byte(tfvars))
	_, err = s3client.PutObject(&s3.PutObjectInput{
		Body:                 fileBytes,
		Bucket:               aws.String(secretsBucket),
		Key:                  aws.String("secrets.tfvars"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
