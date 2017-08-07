package aws

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func hasTable(dynamodbClient *dynamodb.DynamoDB, tableName string) (bool, error) {
	tables, err := dynamodbClient.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return false, err
	}

	for _, table := range tables.TableNames {
		if tableName == *table {
			return true, err
		}
	}
	return false, err
}
