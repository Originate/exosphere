package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

// InitAccount prepares a blank AWS account to be used with Terraform
func InitAccount(options Options) error {
	config := CreateAwsConfig(options)
	session := session.Must(session.NewSession())
	err := createRemoteState(session, config, options.BucketName)
	if err != nil {
		return err
	}
	return createLockTable(session, config, options.TerraformLockTable)
}

// creates s3 bucket to store terraform remote state
func createRemoteState(currSession client.ConfigProvider, config *aws.Config, bucketName string) error {
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
func createLockTable(currSession client.ConfigProvider, config *aws.Config, tableName string) error {
	dynamodbClient := dynamodb.New(currSession, config)
	hasTable, err := hasTable(dynamodbClient, tableName)
	if err != nil {
		return err
	}
	if hasTable {
		return nil
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
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}
	_, err = dynamodbClient.CreateTable(input)
	return err
}
