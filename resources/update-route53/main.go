package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {
	if len(os.Arg) < 3 {
		panic(errors.New("Not enough arguments. Arguments must be non-empty strings passed into 'route53-updater <service-role> <internal-hosted-zone-name>'"))
	}
	serviceRole := os.Args[1]
	internalHostedZoneName := os.Args[2]
	if serviceRole == "" || internalHostedZoneName == "" {
		panic(errors.New("Service role or internal hosted zone name missing. Both arguments must be non-empty strings passed into 'route53-updater <service-role> <internal-hosted-zone-name>'"))
	}
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
	var hostedZoneFound bool
	for _, hostedZone := range hostedZones.HostedZones {
		if *hostedZone.Name == fmt.Sprintf("%s.", internalHostedZoneName) {
			hostedZoneId = hostedZone.Id
			hostedZoneFound = true
		}
	}
	if !hostedZoneFound {
		panic(fmt.Errorf("Hosted zone name '%s' not found.", internalHostedZoneName))
	}

	changeRecordResult, err := awsClient.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: hostedZoneId,
		ChangeBatch: &route53.ChangeBatch{
			Comment: aws.String("Update dns record"),
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(fmt.Sprintf("%s.%s", serviceRole, internalHostedZoneName)),
						Type: aws.String("A"),
						TTL:  aws.Int64(int64(300)),
						ResourceRecords: []*route53.ResourceRecord{
							{Value: aws.String(string(internalIP))},
						},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
}
