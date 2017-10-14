package aws

import (
	"fmt"

	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/moby/moby/client"
)

// RepositoryHelper is used for help pushing an image to ECR
type RepositoryHelper struct {
	EcrAuth        string
	EcrClient      *ecr.ECR
	ImageName      string
	RepositoryName string
	RepositoryURI  string
	ImageVersion   string
}

// GetTaggedImageName returns the image name on ECR
func (r *RepositoryHelper) GetTaggedImageName() string {
	return fmt.Sprintf("%s:%s", r.RepositoryURI, r.ImageVersion)
}

// NeedsPush returns whether or not the image needs to be pushed
func (r *RepositoryHelper) NeedsPush() (bool, error) {
	hasImageVersion, err := r.hasImageVersion()
	if err != nil {
		return false, err
	}
	return hasImageVersion, nil
}

// Push pushes the image to ECR
func (r *RepositoryHelper) Push() error {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	return tools.PushImage(dockerClient, r.GetTaggedImageName(), r.EcrAuth)
}

// Helpers

func (r *RepositoryHelper) hasImageVersion() (bool, error) {
	result, err := r.EcrClient.DescribeImages(&ecr.DescribeImagesInput{
		RepositoryName: aws.String(r.RepositoryName),
	})
	if err != nil {
		return false, err
	}
	for _, imageDetail := range result.ImageDetails {
		for _, tag := range imageDetail.ImageTags {
			if *tag == r.ImageVersion {
				return true, nil
			}
		}
	}
	return false, nil
}
