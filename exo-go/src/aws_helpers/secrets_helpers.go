package awsHelper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const secretsFile string = "secrets.tfvars"

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

func ReadSecrets(secretsBucket, region string) (string, error) {
	s3client := createS3client(region)
	err := verifyBucket(s3client, secretsBucket)
	if err != nil {
		return "", err
	}

	results, err := s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(secretsBucket),
		Key:    aws.String(secretsFile),
	})

	defer results.Body.Close()
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, results.Body); err != nil {
		return "", err
	}

	return string(buffer.Bytes()), err
}

// CreateSecrets creates new secret key value pair
func CreateSecrets(newSecrets map[string]string, secretsBucket, region string) error {
	//TODO: a=""b"" <<this
	tfvars, err := ReadSecrets(secretsBucket, region)
	if err != nil {
		return err
	}

	existingSecrets := tfstringToMap(tfvars)
	if hasConflictingKey(existingSecrets, newSecrets) {
		return errors.New("New secrets have key(s) that conflict with existing secrets. Use 'exo configure update' to update existing keys.")
	}

	return putSecretsFile(mergeMaps(newSecrets, existingSecrets), secretsBucket, region)
}

// takes a map from {"secret_name": "secret_value"} and
// creates a .tfvars file and uploads it to s3
func putSecretsFile(secrets map[string]string, secretsBucket, region string) error {
	s3client := createS3client(region)
	err := verifyBucket(s3client, secretsBucket)
	if err != nil {
		return err
	}

	tfvars := mapToTfstring(secrets)
	fileBytes := bytes.NewReader([]byte(tfvars))

	_, err = s3client.PutObject(&s3.PutObjectInput{
		Body:                 fileBytes,
		Bucket:               aws.String(secretsBucket),
		Key:                  aws.String(secretsFile),
		ServerSideEncryption: aws.String("AES256"),
	})

	return err
}

func mapToTfstring(secrets map[string]string) string {
	tfvars := ""
	for key, value := range secrets {
		tfvars += fmt.Sprintf("%s=\"%s\"\n", key, value)
	}
	return tfvars
}

func tfstringToMap(tfvars string) map[string]string {
	secretsMap := map[string]string{}
	secretPairs := strings.Split(tfvars, "\n")
	secretPairs = secretPairs[:len(secretPairs)-1] //remove trailing empty elem
	for _, secret := range secretPairs {
		s := strings.Split(secret, "=")
		secretsMap[s[0]] = s[1]
	}
	return secretsMap
}

func hasConflictingKey(map1 map[string]string, map2 map[string]string) bool {
	for k, _ := range map2 {
		if _, hasKey := map1[k]; hasKey {
			return true
		}
	}
	return false
}

func mergeMaps(map1 map[string]string, map2 map[string]string) map[string]string {
	for k, v := range map2 {
		map1[k] = v
	}
	return map1
}

// checks if bucket exists and returns a proper error message otherwise
// different from hasBucket which only checks if a bucket exists
func verifyBucket(s3client *s3.S3, secretsBucket string) error {
	hasBucket, err := hasBucket(s3client, secretsBucket)
	if err != nil {
		return err
	}
	if !hasBucket {
		return errors.New("Secrets bucket not found. Run 'exo configure' to initialize bucket.")
	}

	return nil
}
