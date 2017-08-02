package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

// InitAccount prepares a blank AWS account to be used with Terraform
func InitAccount(bucketName, tableName, region string) error {
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
	if hasBucket {
		return nil
	}
	_, err = s3client.CreateBucket(&s3.CreateBucketInput{Bucket: &bucketName})
	if err != nil {
		return err
	}
	_, err = s3client.PutBucketVersioning(&s3.PutBucketVersioningInput{
		Bucket: &bucketName,
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String("Enabled"),
		},
	})
	return err
}

// creates dynamodb table to store terraform lock file
func createLockTable(currSession *session.Session, config *aws.Config, tableName string) error {
	dynamodbClient := dynamodb.New(currSession, config)
	hasTable, err := hasTable(dynamodbClient, tableName)
	if err != nil {
		return err
	}
	if hasTable {
		return err
	}
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
	return err
}
