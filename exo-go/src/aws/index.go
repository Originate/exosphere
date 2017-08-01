package aws

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

// General helpers for AWS operations

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

func createS3client(region string) *s3.S3 {
	config := aws.NewConfig().WithRegion(region)
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
	_, err := s3client.PutObject(&s3.PutObjectInput{
		Body:                 fileContents,
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileName),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
