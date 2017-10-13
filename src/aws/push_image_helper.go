package aws

import (
	"fmt"

	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/moby/moby/client"
)

// PushImageHelper is used for help pushing an image to ECR
type PushImageHelper struct {
	EcrAuth        string
	EcrClient      *ecr.ECR
	ImageName      string
	RepositoryName string
	RepositoryURI  string
	ImageVersion   string
}

// GetTaggedImageName returns the image name on ECR
func (p *PushImageHelper) GetTaggedImageName() string {
	return fmt.Sprintf("%s:%s", p.RepositoryURI, p.ImageVersion)
}

// NeedsPush returns whether or not the image needs to be pushed
func (p *PushImageHelper) NeedsPush() (bool, error) {
	hasImageVersion, err := p.hasImageVersion()
	if err != nil {
		return false, err
	}
	return hasImageVersion, nil
}

// Push pushes the image to ECR
func (p *PushImageHelper) Push() error {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	return tools.PushImage(dockerClient, p.GetTaggedImageName(), p.EcrAuth)
}

// Helpers

func (p *PushImageHelper) hasImageVersion() (bool, error) {
	result, err := p.EcrClient.DescribeImages(&ecr.DescribeImagesInput{
		RepositoryName: aws.String(p.RepositoryName),
	})
	if err != nil {
		return false, err
	}
	for _, imageDetail := range result.ImageDetails {
		for _, tag := range imageDetail.ImageTags {
			if *tag == p.ImageVersion {
				return true, nil
			}
		}
	}
	return false, nil
}
