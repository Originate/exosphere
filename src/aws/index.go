package aws

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/s3"
)

// General helper for AWS operations
func hasBucket(s3client *s3.S3, bucketName string) (bool, error) {
	buckets, err := s3client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return false, err
	}
	for _, bucket := range buckets.Buckets {
		if bucketName == *bucket.Name {
			return true, err
		}
	}
	return false, err
}

// create s3 bucket if it doesn't already exist
func createBucket(s3client *s3.S3, bucketName string) error {
	hasBucket, err := hasBucket(s3client, bucketName)
	if err != nil {
		return err
	}
	if hasBucket {
		return nil
	}
	_, err = s3client.CreateBucket(&s3.CreateBucketInput{Bucket: &bucketName})
	return err
}

func hasTable(dynamodbClient *dynamodb.DynamoDB, tableName string) (bool, error) {
	tables, err := dynamodbClient.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return false, err
	}
	for _, table := range tables.TableNames {
		if tableName == *table {
			return true, nil
		}
	}
	return false, nil
}

func createS3client(awsConfig types.AwsConfig) *s3.S3 {
	config := createAwsConfig(awsConfig)
	currSession := session.Must(session.NewSession())
	return s3.New(currSession, config)
}

// create s3 bucket and object if it doesn't already exist
func createS3Object(s3client *s3.S3, fileContents io.ReadSeeker, bucketName, fileName string) error {
	err := createBucket(s3client, bucketName)
	if err != nil {
		return err
	}
	_, err = s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return putS3Object(s3client, fileContents, bucketName, fileName)
		}
		return err
	}
	return nil
}

func putS3Object(s3client *s3.S3, fileContents io.ReadSeeker, bucketName, fileName string) error {
	err := createBucket(s3client, bucketName)
	if err != nil {
		return err
	}
	_, err = s3client.PutObject(&s3.PutObjectInput{
		Body:                 fileContents,
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileName),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

// retrieves repository URI given a repository name
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

// creates an image repository if it doesn't already exist and returns its URI
func createRepository(ecrClient *ecr.ECR, repositoryName string) (string, error) {
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

// retrieves encoded ECR credentails (in the format username:password) and returns them as separate strings
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

func createAwsConfig(awsConfig types.AwsConfig) *aws.Config {
	return &aws.Config{
		Region: aws.String(awsConfig.Region),
		Credentials: credentials.NewCredentials(&credentials.ChainProvider{
			Providers: []credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{Profile: awsConfig.Profile},
			},
		}),
	}
}
