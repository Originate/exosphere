package terraformFileHelpers_test

import (
	"os"
	"regexp"

	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/terraform_file_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
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
	homeDir, err = osHelpers.GetUserHomeDir()
	if err != nil {
		panic(err)
	}
})

var _ = Describe("Given an application with no services", func() {
	appConfig := types.AppConfig{
		Name: "example-app",
	}
	serviceConfigs := map[string]types.ServiceConfig{}

	It("should generate an AWS module only", func() {
		result, err := terraformFileHelpers.GenerateTerraform(appConfig, serviceConfigs, appDir, homeDir)
		Expect(err).To(BeNil())
		expected := normalizeWhitespace(
			`terraform {
	required_version = "= 0.9.9"

	backend "s3" {
		bucket     = "example-app-terraform"
		key        = "dev/terraform.tfstate"
		region     = "us-west-2"
		lock_table = "TerraformLocks"
	}
}

module "aws" {
	source = "./aws"

	account_id       = "${var.account_id}"
	application_name = "example-app"
	env              = "production"
	key_name         = "${var.key_name}"
	region           = "${var.region}"
	security_groups  = []
}`)
		Expect(result).To(ContainSubstring(expected))
	})
})

var _ = Describe("Given an application with public and private services", func() {
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
	result, err := terraformFileHelpers.GenerateTerraform(appConfig, serviceConfigs, appDir, homeDir)

	It("should generate a public service module", func() {
		Expect(err).To(BeNil())
		expected := normalizeWhitespace(
			`module "public-service" {
  source = "./aws/public-service"

  name = "public-service"

  alb_security_group      = ["${module.aws.external_alb_security_group}"]
  alb_subnet_ids          = ["${module.aws.public_subnet_ids}"]
  cluster_id              = "${module.aws.cluster_id}"
  command                 = ["node","app"]
  container_port          = "3000"
  cpu_units               = "128"
  ecs_role_arn            = "${module.aws.ecs_service_iam_role_arn}"
  env                     = "production"
  external_dns_name       = "originate.com"
  external_hosted_zone_id = "${var.hosted_zone_id}"
  health_check_endpoint   = "/health-check"
  internal_dns_name       = "${module.aws.internal_dns_name}"
  internal_hosted_zone_id = "${module.aws.internal_hosted_zone_id}"
  log_bucket              = "${module.aws.log_bucket_id}"
  memory_reservation      = "128"
  region                  = "${var.region}"
  vpc_id                  = "${module.aws.vpc_id}"
}`)
		Expect(result).To(ContainSubstring(expected))
	})

	It("should generate a private service module", func() {
		Expect(err).To(BeNil())
		expected := normalizeWhitespace(
			`module "private-service" {
  source = "./aws/worker-service"

  name = "private-service"

  cluster_id            = "${module.aws.cluster_id}"
  command               = ["exo-js"]
  cpu_units             = "128"
  env                   = "production"
  memory_reservation    = "128"
  region                = "${var.region}"
}`)
		Expect(result).To(ContainSubstring(expected))
	})
})

var _ = Describe("Given an application with dependencies", func() {
	dependencies := []types.Dependency{
		types.Dependency{
			Name:    "exocom",
			Version: "0.22.1",
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

	It("should generate dependency modules", func() {
		result, err := terraformFileHelpers.GenerateTerraform(appConfig, serviceConfigs, appDir, homeDir)
		Expect(err).To(BeNil())
		expected := normalizeWhitespace(
			`module "exocom_cluster" {
  source = "./aws/custom/exocom/exocom-cluster"

  availability_zones = "${module.aws.availability_zones}"
  env                = "production"
  domain_name        = "originate.com"
  hosted_zone_id     = "${var.hosted_zone_id}"
  instance_type      = "t2.micro"
  key_name           = "${var.key_name}"
  name               = "exocom"
  region             = "${module.aws.region}"
  security_groups = ["${module.aws.bastion_security_group}",
    "${module.aws.ecs_cluster_security_group}",
    "${module.aws.external_alb_security_group}",
  ]
  subnet_ids = "${module.aws.private_subnet_ids}"
  vpc_id     = "${module.aws.vpc_id}"
}

module "exocom_service" {
  source = "./aws/custom/exocom/exocom-service"

  cluster_id     = "${module.exocom_cluster.cluster_id}"
  command        = ["bin/exocom"]
  container_port = "3100"
  cpu_units      = "128"
  env            = "production"
  environment_variables = {
    ROLE  = "exocom"
  }
  memory_reservation = "128"
  name               = "exocom"
  region             = "${module.aws.region}"
}`)
		Expect(result).To(ContainSubstring(expected))
	})
})

func normalizeWhitespace(str string) string {
	regex := regexp.MustCompile("\t")
	return regex.ReplaceAllString(str, "  ")
}
