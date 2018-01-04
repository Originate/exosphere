package main

import (
	"fmt"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {
	// config := &aws.Config{
	// 	Region: aws.String("us-west-2"),
	// 	Credentials: credentials.NewCredentials(&credentials.ChainProvider{
	// 		Providers: []credentials.Provider{
	// 			&credentials.EnvProvider{},
	// 			&credentials.SharedCredentialsProvider{Profile: "space-tweet"},
	// 		},
	// 	}),
	// }
	internalIP, err := exec.Command("curl", "-fsSL", "http://169.254.169.254/latest/meta-data/local-ipv4").Output()
	if err != nil {
		panic(err)
	}

	awsClient := route53.New(session.New())
	hostedZones, err := awsClient.ListHostedZones(&route53.ListHostedZonesInput{})
	if err != nil {
		panic(err)
	}
	var hostedZoneId *string
	for _, hostedZone := range hostedZones.HostedZones {
		if *hostedZone.Name == "space-tweet.local." {
			hostedZoneId = hostedZone.Id
		}
	}
	changeRecordResult, err := awsClient.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: hostedZoneId,
		ChangeBatch: &route53.ChangeBatch{
			Comment: aws.String("Update dns record"),
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String("test.space-tweet.local"),
						Type: aws.String("A"),
						TTL:  aws.Int64(int64(300)),
						ResourceRecords: []*route53.ResourceRecord{
							{Value: aws.String(internalIP)},
						},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(changeRecordResult)
}
