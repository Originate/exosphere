package terraform_test

import (
	"encoding/json"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/Originate/exosphere/test_helpers"
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

var _ = Describe("Terraform commands", func() {
	service1EnvVars := types.EnvVars{
		Default: map[string]string{
			"env1": "val1",
		},
		Secrets: []string{"secret1"},
	}
	service1Config := types.ServiceConfig{
		Environment: service1EnvVars,
	}
	secrets := map[string]string{
		"secret1": "secret_value1",
	}
	deployConfig := types.DeployConfig{
		ServiceConfigs: map[string]types.ServiceConfig{
			"service1": service1Config,
		},
	}

	It("should compile the proper var flags", func() {
		vars, err := terraform.CompileVarFlags(deployConfig, secrets)
		Expect(err).NotTo(HaveOccurred())
		Expect(vars[0]).To(Equal("-var"))
		Expect(vars[1]).To(Equal("secret1=secret_value1"))
		Expect(vars[2]).To(Equal("-var"))

		varName := strings.Split(vars[3], "=")[0]
		varVal := strings.Split(vars[3], "=")[1]
		var escapedVal string
		actualVal := []map[string]string{}
		expectedVal := []map[string]string{
			{
				"key":   "secret1",
				"value": "secret_value1",
			},
			{
				"key":   "env1",
				"value": "val1",
			},
		}
		err = json.Unmarshal([]byte(varVal), &escapedVal)
		Expect(err).NotTo(HaveOccurred())
		err = json.Unmarshal([]byte(escapedVal), &actualVal)
		Expect(err).NotTo(HaveOccurred())
		Expect(varName).To(Equal("service1_env_vars"))
		Expect(expectedVal).To(ConsistOf(actualVal))
	})
})

var _ = Describe("Template builder", func() {
	var _ = Describe("Given an application with no services", func() {
		appConfig := types.AppConfig{
			Name: "example-app",
			Production: map[string]string{
				"url": "example-app.com",
			},
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
			result, err := terraform.Generate(deployConfig, map[string]string{})
			Expect(err).To(BeNil())
			expected := normalizeWhitespace(
				`terraform {
	required_version = "= 0.10.0"

	backend "s3" {
		bucket         = "example-app-terraform"
		key            = "dev/terraform.tfstate"
		region         = "us-west-2"
		dynamodb_table = "TerraformLocks"
	}
}

provider "aws" {
  version = "0.1.4"

  region              = "${var.region}"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["${var.account_id}"]
}

module "aws" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws?ref=8786f912"

  name              = "example-app"
  env               = "production"
	external_dns_name = "example-app.com"
  key_name          = "${var.key_name}"
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
		imagesMap := map[string]string{
			"public-service":  "test-public-image:0.0.1",
			"private-service": "test-private-image:0.0.1",
		}

		BeforeEach(func() {
			var err error
			result, err = terraform.Generate(deployConfig, imagesMap)
			Expect(err).To(BeNil())
		})

		It("should generate a public service module", func() {
			expected := normalizeWhitespace(
				`variable "public-service_env_vars" {
  default = "[]"
}

module "public-service" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//public-service?ref=8786f912"

  name = "public-service"

  alb_security_group    = "${module.aws.external_alb_security_group}"
  alb_subnet_ids        = ["${module.aws.public_subnet_ids}"]
  cluster_id            = "${module.aws.ecs_cluster_id}"
  command               = ["node","app"]
  container_port        = "3000"
  cpu                   = "128"
  desired_count         = 1
	docker_image          = "test-public-image:0.0.1"
  ecs_role_arn          = "${module.aws.ecs_service_iam_role_arn}"
  env                   = "production"
  environment_variables = "${var.public-service_env_vars}"
  external_dns_name     = "originate.com"
  external_zone_id      = "${module.aws.external_zone_id}"
  health_check_endpoint = "/health-check"
  internal_dns_name     = "public-service"
  internal_zone_id      = "${module.aws.internal_zone_id}"
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
				`variable "private-service_env_vars" {
  default = "[]"
}

module "private-service" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//worker-service?ref=8786f912"

  name = "private-service"

  cluster_id            = "${module.aws.ecs_cluster_id}"
  command               = ["exo-js"]
  cpu                   = "128"
  desired_count         = 1
	docker_image          = "test-private-image:0.0.1"
  env                   = "production"
  environment_variables = "${var.private-service_env_vars}"
  memory                = "128"
  region                = "${var.region}"
}`)
			Expect(result).To(ContainSubstring(expected))
		})
	})

	var _ = Describe("Given an application with dependencies", func() {

		It("should generate dependency modules", func() {
			cwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			err = testHelpers.CheckoutApp(cwd, "simple")
			Expect(err).NotTo(HaveOccurred())
			appDir := path.Join("tmp", "simple")
			homeDir, err := util.GetHomeDirectory()
			if err != nil {
				panic(err)
			}
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
			Expect(err).NotTo(HaveOccurred())

			deployConfig := types.DeployConfig{
				AppConfig:      appConfig,
				ServiceConfigs: serviceConfigs,
				AppDir:         appDir,
				HomeDir:        homeDir,
			}
			imagesMap := map[string]string{
				"exocom": "originate/exocom:0.0.1",
			}

			result, err := terraform.Generate(deployConfig, imagesMap)
			Expect(err).To(BeNil())
			expected := normalizeWhitespace(
				`module "exocom_cluster" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//custom//exocom//exocom-cluster?ref=8786f912"

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
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//custom//exocom//exocom-service?ref=8786f912"

  cluster_id            = "${module.exocom_cluster.cluster_id}"
  container_port        = "3100"
  cpu_units             = "128"
	docker_image          = "originate/exocom:0.0.1"
  env                   = "production"
  environment_variables = {
    ROLE = "exocom"
		SERVICE_ROUTES = <<EOF
[{"receives":["users.created"],"role":"web","sends":["users.create"]}]
EOF
  }
  memory_reservation    = "128"
  name                  = "exocom"
  region                = "${module.aws.region}"
}`)
			Expect(result).To(ContainSubstring(expected))
		})
	})
})

func normalizeWhitespace(str string) string {
	regex := regexp.MustCompile("\t")
	return regex.ReplaceAllString(str, "  ")
}
