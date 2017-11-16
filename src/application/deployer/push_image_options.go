package deployer

import (
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// PushImageOptions is the options to PushServiceImage
type PushImageOptions struct {
	DeployConfig    types.DeployConfig
	EcrAuth         string
	EcrClient       *ecr.ECR
	ImageName       string
	ServiceLocation string
	ServiceRole     string
	BuildImage      bool
	BuildMode       composebuilder.BuildMode
}
