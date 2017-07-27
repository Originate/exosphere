package awsHelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Configure sets up a blank AWS account to have s3 bucket and dynamodb table
func Configure(bucketName, tableName, region string) error {
	config := aws.NewConfig().WithRegion(region)
	session := session.Must(session.NewSession())

	err := createRemoteState(session, config, bucketName)
	if err != nil {
		return err
	}

	return createLockTable(session, config, tableName)
}

// creates s3 bucket to store terraform remote state
func createRemoteState(currSession *session.Session, config *aws.Config, bucketName string) error {
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

		// enable versioning
		versioningStatus := "Enabled"
		_, err = s3client.PutBucketVersioning(&s3.PutBucketVersioningInput{
			Bucket: &bucketName,
			VersioningConfiguration: &s3.VersioningConfiguration{
				Status: &versioningStatus,
			},
		})
	}

	return err
}

// creates dynamodb table to store terraform lock file
func createLockTable(currSession *session.Session, config *aws.Config, tableName string) error {
	dynamodbClient := dynamodb.New(currSession, config)

	hasTable, err := hasTable(dynamodbClient, tableName)

	if !hasTable {
		input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("LockID"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("LockID"),
					KeyType:       aws.String("HASH"),
				},
			},
			TableName: aws.String(tableName),
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		}

		_, err = dynamodbClient.CreateTable(input)
	}

	return err
}
