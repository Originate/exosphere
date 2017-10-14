package deployer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/types"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// PushServiceImageOptions is the options to PushServiceImage
type PushServiceImageOptions struct {
	DeployConfig     types.DeployConfig
	DockerComposeDir string
	EcrAuth          string
	EcrClient        *ecr.ECR
	ImageName        string
	ServiceLocation  string
	BuildImage       bool
}

// PushServiceImage pushes a single service image to ECR, building or pulling if needed
func PushServiceImage(options PushServiceImageOptions) (string, error) {
	repositoryHelper, err := getRepositoryHelper(options)
	if err != nil {
		return "", err
	}
	needsPush, err := repositoryHelper.NeedsPush()
	if err != nil {
		return "", err
	}
	if needsPush {
		err := buildOrPullImage(options)
		if err != nil {
			return "", err
		}
		options.DeployConfig.Logger.Logf("Pushing image: %s...", options.ImageName)
		err = repositoryHelper.Push()
		if err != nil {
			return "", err
		}
	} else {
		options.DeployConfig.Logger.Logf("Image %s is up to date, skipping...", options.ImageName)
	}
	return repositoryHelper.GetTaggedImageName(), nil
}

func buildOrPullImage(options PushServiceImageOptions) error {
	opts := compose.ImageOptions{
		DockerComposeDir: options.DockerComposeDir,
		Logger:           options.DeployConfig.Logger,
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DeployConfig.DockerComposeProjectName)},
		ImageName:        options.ImageName,
	}
	if options.ServiceLocation != "" {
		options.DeployConfig.Logger.Logf("Building image: %s...", options.ImageName)
		return compose.BuildImage(opts)
	}
	options.DeployConfig.Logger.Logf("Pulling image: %s...", options.ImageName)
	return compose.PullImage(opts)
}

func getRepositoryHelper(options PushServiceImageOptions) (*aws.RepositoryHelper, error) {
	imageNameParts := strings.Split(options.ImageName, ":")
	repositoryName := imageNameParts[0]
	var imageVersion string
	if options.ServiceLocation != "" {
		var err error
		imageVersion, err = getCommitSHA(options.DeployConfig.AppDir, options.ServiceLocation)
		if err != nil {
			return nil, err
		}
	} else {
		imageVersion = imageNameParts[1]
	}
	repositoryURI, err := aws.CreateRepository(options.EcrClient, repositoryName)
	if err != nil {
		return nil, err
	}
	return &aws.RepositoryHelper{
		EcrAuth:        options.EcrAuth,
		EcrClient:      options.EcrClient,
		ImageName:      options.ImageName,
		ImageVersion:   imageVersion,
		RepositoryName: repositoryName,
		RepositoryURI:  repositoryURI,
	}, nil
}

func getCommitSHA(appDir string, serviceLocation string) (string, error) {
	cmd := exec.Command("git", "rev-list", "-1", "HEAD", serviceLocation)
	cmd.Dir = appDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(output), "\n"), err
}
