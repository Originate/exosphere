package terraform_test

import (
	"os"
	"path"
	"regexp"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/Originate/exosphere/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template builder", func() {
	var _ = Describe("Given an application with no services", func() {
		appConfig := types.AppConfig{
			Name: "example-app",
			Production: types.AppProductionConfig{
				URL: "example-app.com",
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
				AccountID:            "12345",
			},
			TerraformModulesRef: "TERRAFORM_MODULES_REF",
		}

		It("should generate an AWS module only", func() {
			result, err := terraform.Generate(deployConfig, map[string]string{})
			Expect(err).To(BeNil())
			expected := normalizeWhitespace(
				`terraform {
	required_version = ">= 0.10.0"

	backend "s3" {
		bucket         = "example-app-terraform"
		key            = "terraform.tfstate"
		region         = "us-west-2"
		dynamodb_table = "TerraformLocks"
	}
}

provider "aws" {
  version = "0.1.4"

  region              = "us-west-2"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["12345"]
}

variable "key_name" {
  default = ""
}

module "aws" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws?ref=TERRAFORM_MODULES_REF"

  name              = "example-app"
  env               = "production"
	external_dns_name = "example-app.com"
  key_name          = "${var.key_name}"
}`)
			Expect(result).To(ContainSubstring(expected))
		})
	})

	var _ = Describe("Given an application with public, private, and worker services", func() {
		var result string
		services := types.Services{
			Public:  map[string]types.ServiceData{"public-service": {}},
			Private: map[string]types.ServiceData{"private-service": {}},
			Worker:  map[string]types.ServiceData{"worker-service": {}},
		}
		appConfig := types.AppConfig{
			Name:     "example-app",
			Services: services,
		}
		serviceConfigs := map[string]types.ServiceConfig{
			"public-service": {
				Production: types.ServiceProductionConfig{
					PublicPort:  "3000",
					CPU:         "128",
					URL:         "originate.com",
					HealthCheck: "/health-check",
					Memory:      "128",
				},
			},
			"private-service": {
				Production: types.ServiceProductionConfig{
					PublicPort:  "3100",
					CPU:         "128",
					HealthCheck: "/health-check",
					Memory:      "128",
				},
			},
			"worker-service": {
				Production: types.ServiceProductionConfig{
					CPU:    "128",
					Memory: "128",
				},
			},
		}

		deployConfig := types.DeployConfig{
			AppConfig:      appConfig,
			ServiceConfigs: serviceConfigs,
			AppDir:         appDir,
			HomeDir:        homeDir,
			AwsConfig: types.AwsConfig{
				SslCertificateArn: "sslcert123",
			},
			TerraformModulesRef: "TERRAFORM_MODULES_REF",
		}

		BeforeEach(func() {
			var err error
			result, err = terraform.Generate(deployConfig, map[string]string{})
			Expect(err).To(BeNil())
		})

		It("should generate a public service module", func() {
			expected := normalizeWhitespace(
				`variable "public-service_env_vars" {
  default = "[]"
}

variable "public-service_docker_image" {}

module "public-service" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//public-service?ref=TERRAFORM_MODULES_REF"

  name = "public-service"

  alb_security_group    = "${module.aws.external_alb_security_group}"
  alb_subnet_ids        = ["${module.aws.public_subnet_ids}"]
  cluster_id            = "${module.aws.ecs_cluster_id}"
  container_port        = "3000"
  cpu                   = "128"
  desired_count         = 1
	docker_image          = "${var.public-service_docker_image}"
  ecs_role_arn          = "${module.aws.ecs_service_iam_role_arn}"
  env                   = "production"
  environment_variables = "${var.public-service_env_vars}"
  external_dns_name     = "originate.com"
  external_zone_id      = "${module.aws.external_zone_id}"
  health_check_endpoint = "/health-check"
  internal_dns_name     = "public-service"
  internal_zone_id      = "${module.aws.internal_zone_id}"
  log_bucket            = "${module.aws.log_bucket_id}"
  memory_reservation    = "128"
  region                = "${module.aws.region}"
  ssl_certificate_arn   = "sslcert123"
  vpc_id                = "${module.aws.vpc_id}"
}`)
			Expect(result).To(ContainSubstring(expected))
		})

		It("should generate a private service module", func() {
			expected := normalizeWhitespace(
				`variable "private-service_env_vars" {
  default = "[]"
}

variable "private-service_docker_image" {}

module "private-service" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//private-service?ref=TERRAFORM_MODULES_REF"

  name = "private-service"

  alb_security_group    = "${module.aws.internal_alb_security_group}"
  alb_subnet_ids        = ["${module.aws.private_subnet_ids}"]
  cluster_id            = "${module.aws.ecs_cluster_id}"
  container_port        = "3100"
  cpu                   = "128"
  desired_count         = 1
	docker_image          = "${var.private-service_docker_image}"
  ecs_role_arn          = "${module.aws.ecs_service_iam_role_arn}"
  env                   = "production"
  environment_variables = "${var.private-service_env_vars}"
  health_check_endpoint = "/health-check"
  internal_dns_name     = "private-service"
  internal_zone_id      = "${module.aws.internal_zone_id}"
  log_bucket            = "${module.aws.log_bucket_id}"
  memory_reservation    = "128"
  region                = "${module.aws.region}"
  vpc_id                = "${module.aws.vpc_id}"
}`)
			Expect(result).To(ContainSubstring(expected))
		})

		It("should generate a worker service module", func() {
			expected := normalizeWhitespace(
				`variable "worker-service_env_vars" {
  default = "[]"
}

variable "worker-service_docker_image" {}

module "worker-service" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//worker-service?ref=TERRAFORM_MODULES_REF"

  name = "worker-service"

  cluster_id            = "${module.aws.ecs_cluster_id}"
  cpu                   = "128"
  desired_count         = 1
	docker_image          = "${var.worker-service_docker_image}"
  env                   = "production"
  environment_variables = "${var.worker-service_env_vars}"
  memory_reservation    = "128"
  region                = "${module.aws.region}"
}`)
			Expect(result).To(ContainSubstring(expected))
		})
	})

	var _ = Describe("Given an application with dependencies", func() {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		homeDir, err := util.GetHomeDirectory()
		if err != nil {
			panic(err)
		}

		It("should generate dependency modules for exocom", func() {
			err = testHelpers.CheckoutApp(cwd, "simple")
			Expect(err).NotTo(HaveOccurred())
			appDir := path.Join("tmp", "simple")
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
			Expect(err).NotTo(HaveOccurred())

			deployConfig := types.DeployConfig{
				AppConfig:           appConfig,
				ServiceConfigs:      serviceConfigs,
				AppDir:              appDir,
				HomeDir:             homeDir,
				TerraformModulesRef: "TERRAFORM_MODULES_REF",
			}
			imagesMap := map[string]string{
				"exocom": "originate/exocom:0.0.1",
			}

			result, err := terraform.Generate(deployConfig, imagesMap)
			Expect(err).To(BeNil())
			expected := normalizeWhitespace(
				`module "exocom_cluster" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//dependencies//exocom//exocom-cluster?ref=TERRAFORM_MODULES_REF"

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

variable "exocom_service_routes" {
	default = ""
}

module "exocom_service" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//dependencies//exocom//exocom-service?ref=TERRAFORM_MODULES_REF"

  cluster_id            = "${module.exocom_cluster.cluster_id}"
  cpu_units             = "128"
	docker_image          = "originate/exocom:0.0.1"
  env                   = "production"
  environment_variables = {
		SERVICE_ROUTES = "${var.exocom_service_routes}"
  }
  memory_reservation    = "128"
  name                  = "exocom"
  region                = "${module.aws.region}"
}`)
			Expect(result).To(ContainSubstring(expected))
		})

		It("should generate rds modules for dependencies", func() {
			err = testHelpers.CheckoutApp(cwd, "rds")
			Expect(err).NotTo(HaveOccurred())
			appDir := path.Join("tmp", "rds")
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
			Expect(err).NotTo(HaveOccurred())

			deployConfig := types.DeployConfig{
				AppConfig:           appConfig,
				AppDir:              appDir,
				ServiceConfigs:      serviceConfigs,
				TerraformModulesRef: "TERRAFORM_MODULES_REF",
			}
			imagesMap := map[string]string{
				"postgres": "postgres:9.6.4",
				"mysql":    "mysql:5.6.17",
			}

			result, err := terraform.Generate(deployConfig, imagesMap)
			Expect(err).To(BeNil())
			By("generating rds modules for application dependencies", func() {
				expected := normalizeWhitespace(
					`module "my-db_rds_instance" {
	source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//dependencies//rds?ref=TERRAFORM_MODULES_REF"

  allocated_storage       = 10
  ecs_security_group      = "${module.aws.ecs_cluster_security_group}"
  engine                  = "postgres"
  engine_version          = "9.6.4"
  env                     = "production"
  instance_class          = "db.t2.micro"
  internal_hosted_zone_id = "${module.aws.internal_zone_id}"
  name                    = "my-db"
  username                = "originate-user"
  password                = "${var.POSTGRES_PASSWORD}"
  storage_type            = "gp2"
  subnet_ids              = ["${module.aws.private_subnet_ids}"]
  vpc_id                  = "${module.aws.vpc_id}"
}`)
				Expect(result).To(ContainSubstring(expected))
			})

			By("should generate rds modules for service dependencies", func() {
				expected := normalizeWhitespace(
					`module "my-sql-db_rds_instance" {
	source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//dependencies//rds?ref=TERRAFORM_MODULES_REF"

  allocated_storage       = 10
  ecs_security_group      = "${module.aws.ecs_cluster_security_group}"
  engine                  = "mysql"
  engine_version          = "5.6.17"
  env                     = "production"
  instance_class          = "db.t1.micro"
  internal_hosted_zone_id = "${module.aws.internal_zone_id}"
  name                    = "my-sql-db"
  username                = "originate-user"
  password                = "${var.MYSQL_PASSWORD}"
  storage_type            = "gp2"
  subnet_ids              = ["${module.aws.private_subnet_ids}"]
  vpc_id                  = "${module.aws.vpc_id}"
}`)
				Expect(result).To(ContainSubstring(expected))
			})
		})
	})
})

func normalizeWhitespace(str string) string {
	regex := regexp.MustCompile("\t")
	return regex.ReplaceAllString(str, "  ")
}
