package deployer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/tools"
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
		err = tools.TagImage(options.ImageName, repositoryHelper.GetTaggedImageName())
		if err != nil {
			return "", err
		}
		fmt.Fprintf(options.DeployConfig.Writer, "Pushing image: %s...\n", options.ImageName)
		err = repositoryHelper.Push(options.DeployConfig.Writer)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Fprintf(options.DeployConfig.Writer, "Image %s is up to date, skipping...\n", options.ImageName)
	}
	return repositoryHelper.GetTaggedImageName(), nil
}

func buildOrPullImage(options PushImageOptions) error {
	opts := compose.CommandOptions{
		DockerComposeDir: options.DockerComposeDir,
		Writer:           options.DeployConfig.Writer,
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DeployConfig.DockerComposeProjectName)},
	}
	if options.ServiceLocation != "" {
		opts.ImageNames = []string{options.ServiceRole}
		fmt.Fprintf(options.DeployConfig.Writer, "Building image: %s...\n", options.ImageName)
		return compose.BuildImages(opts)
	}
	opts.ImageNames = []string{options.ImageName}
	fmt.Fprintf(options.DeployConfig.Writer, "Pulling image: %s...\n", options.ImageName)
	return compose.PullImages(opts)
}

func getRepositoryHelper(options PushImageOptions) (*aws.RepositoryHelper, error) {
	imageNameParts := strings.Split(options.ImageName, ":")
	repositoryName := imageNameParts[0]
	var imageVersion string
	if options.ServiceLocation != "" {
		var err error
		imageVersion, err = getCommitSHA(options.DeployConfig.AppContext.Location, options.ServiceLocation)
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
