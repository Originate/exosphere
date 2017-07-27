package awsHelper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// CreateSecretsStore creates an S3 bucket used for secrets management
func CreateSecretsStore(bucketName, region string) error {
	config := aws.NewConfig().WithRegion(region)
	currSession := session.Must(session.NewSession())
	s3client := s3.New(currSession, config)

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

func CreateSecretEntry(entries []string) error {
	entryMap := map[string]string{}
	for _, entry := range entries {
		s := strings.Split(entry, "=")
		entryMap[s[0]] = s[1]
	}
	secrets, _ := json.Marshal(entryMap)
	fmt.Println(string(secrets))
	return nil
}

func UpdateSecretEntry() error {
	return nil

}

func DeleteSecretEntry() error {
	return nil

}
