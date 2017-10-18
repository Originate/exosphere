package deployer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/compose"
)

// PushImage pushes a single service/dependency image to ECR, building or pulling if needed
func PushImage(options PushImageOptions) (string, error) {
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

func buildOrPullImage(options PushImageOptions) error {
	opts := compose.ImageOptions{
		DockerComposeDir: options.DockerComposeDir,
		Logger:           options.DeployConfig.Logger,
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DeployConfig.DockerComposeProjectName)},
	}
	if options.ServiceLocation != "" {
		opts.ImageName = options.ServiceRole
		options.DeployConfig.Logger.Logf("Building image: %s...", options.ImageName)
		return compose.BuildImage(opts)
	}
	opts.ImageName = options.ImageName
	options.DeployConfig.Logger.Logf("Pulling image: %s...", options.ImageName)
	return compose.PullImage(opts)
}

func getRepositoryHelper(options PushImageOptions) (*aws.RepositoryHelper, error) {
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
