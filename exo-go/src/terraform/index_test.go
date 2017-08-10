package terraform_test

import (
	"os"
	"regexp"

	"github.com/Originate/exosphere/exo-go/src/terraform"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var appDir string
var homeDir string

var _ = BeforeSuite(func() {
	var err error
	appDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	homeDir, err = util.GetHomeDirectory()
	if err != nil {
		panic(err)
	}
})

var _ = Describe("Given an application with no services", func() {
	appConfig := types.AppConfig{
		Name: "example-app",
	}
	serviceConfigs := map[string]types.ServiceConfig{}

	deployConfig := types.DeployConfig{
		AppConfig:      appConfig,
		ServiceConfigs: serviceConfigs,
		AppDir:         appDir,
		HomeDir:        homeDir,
		AwsConfig: types.AwsConfig{
			TerraformStateBucket: "example-app-terraform",
			TerraformLockTable:   "TerraformLocks",
			Region:               "us-west-2",
		},
	}

	It("should generate an AWS module only", func() {
		result, err := terraform.Generate(deployConfig)
		Expect(err).To(BeNil())
		expected := normalizeWhitespace(
			`terraform {
	required_version = "= 0.9.11"

	backend "s3" {
		bucket     = "example-app-terraform"
		key        = "dev/terraform.tfstate"
		region     = "us-west-2"
		lock_table = "TerraformLocks"
	}
}

provider "aws" {
  region              = "${var.region}"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["${var.account_id}"]
}

module "aws" {
  source = "./aws"

  name     = "example-app"
  env      = "production"
  key_name = "${var.key_name}"
}`)
		Expect(result).To(ContainSubstring(expected))
	})
})

var _ = Describe("Given an application with public and private services", func() {
	var result string
	services := types.Services{
		Public:  map[string]types.ServiceData{"public-service": {}},
		Private: map[string]types.ServiceData{"private-service": {}},
	}
	appConfig := types.AppConfig{
		Name:     "example-app",
		Services: services,
	}
	serviceConfigs := map[string]types.ServiceConfig{
		"public-service": {
			Startup: map[string]string{
				"command": "node app",
			},
			Production: map[string]string{
				"public-port":  "3000",
				"cpu":          "128",
				"url":          "originate.com",
				"health-check": "/health-check",
				"memory":       "128",
			},
		},
		"private-service": {
			Startup: map[string]string{
				"command": "exo-js",
			},
			Production: map[string]string{
				"cpu":    "128",
				"memory": "128",
			},
		},
	}

	deployConfig := types.DeployConfig{
		AppConfig:      appConfig,
		ServiceConfigs: serviceConfigs,
		AppDir:         appDir,
		HomeDir:        homeDir,
	}

	BeforeEach(func() {
		var err error
		result, err = terraform.Generate(deployConfig)
		Expect(err).To(BeNil())
	})

	It("should generate a public service module", func() {
		expected := normalizeWhitespace(
			`module "public-service" {
  source = "./aws/public-service"

  name = "public-service"

  alb_security_group    = ["${module.aws.external_alb_security_group}"]
  alb_subnet_ids        = ["${module.aws.public_subnet_ids}"]
  cluster_id            = "${module.aws.cluster_id}"
  command               = ["node","app"]
  container_port        = "3000"
  cpu                   = "128"
  desired_count         = 1
  ecs_role_arn          = "${module.aws.ecs_service_iam_role_arn}"
  env                   = "production"
  external_dns_name     = "originate.com"
  external_zone_id      = "${var.hosted_zone_id}"
  health_check_endpoint = "/health-check"
  internal_dns_name     = "${module.aws.internal_dns_name}"
  internal_zone_id      = "${module.aws.internal_hosted_zone_id}"
  log_bucket            = "${module.aws.log_bucket_id}"
  memory                = "128"
  region                = "${var.region}"
  ssl_certificate_arn   = "${var.ssl_certificate_arn}"
  vpc_id                = "${module.aws.vpc_id}"
}`)
		Expect(result).To(ContainSubstring(expected))
	})

	It("should generate a private service module", func() {
		expected := normalizeWhitespace(
			`module "private-service" {
  source = "./aws/worker-service"

  name = "private-service"

  cluster_id    = "${module.aws.cluster_id}"
  command       = ["exo-js"]
  cpu           = "128"
  desired_count = 1
  env           = "production"
  memory        = "128"
  region        = "${var.region}"
}`)
		Expect(result).To(ContainSubstring(expected))
	})
})

var _ = Describe("Given an application with dependencies", func() {
	dependencies := []types.Dependency{
		types.Dependency{
			Name:    "exocom",
			Version: "0.24.0",
		},
	}
	appConfig := types.AppConfig{
		Name:         "example-app",
		Dependencies: dependencies,
		Production: map[string]string{
			"url": "originate.com",
		},
	}
	serviceConfigs := map[string]types.ServiceConfig{}

	deployConfig := types.DeployConfig{
		AppConfig:      appConfig,
		ServiceConfigs: serviceConfigs,
		AppDir:         appDir,
		HomeDir:        homeDir,
	}

	It("should generate dependency modules", func() {
		result, err := terraform.Generate(deployConfig)
		Expect(err).To(BeNil())
		expected := normalizeWhitespace(
			`module "exocom_cluster" {
  source = "./aws/custom/exocom/exocom-cluster"

  availability_zones          = "${module.aws.availability_zones}"
  env                         = "production"
  internal_hosted_zone_id     = "${module.aws.internal_zone_id}"
  instance_type               = "t2.micro"
  key_name                    = "${var.key_name}"
  name                        = "exocom"
  region                      = "${module.aws.region}"

  bastion_security_group      = ["${module.aws.bastion_security_group}"]

  ecs_cluster_security_groups = [ "${module.aws.ecs_cluster_security_group}",
    "${module.aws.external_alb_security_group}",
  ]

  subnet_ids                  = "${module.aws.private_subnet_ids}"
  vpc_id                      = "${module.aws.vpc_id}"
}

module "exocom_service" {
  source = "./aws/custom/exocom/exocom-service"

  cluster_id            = "${module.exocom_cluster.cluster_id}"
  command               = ["bin/exocom"]
  container_port        = "3100"
  cpu_units             = "128"
  env                   = "production"
  environment_variables = {
    ROLE = "exocom"
  }
  memory_reservation    = "128"
  name                  = "exocom"
  region                = "${module.aws.region}"
}`)
		Expect(result).To(ContainSubstring(expected))
	})
})

func normalizeWhitespace(str string) string {
	regex := regexp.MustCompile("\t")
	return regex.ReplaceAllString(str, "  ")
}
