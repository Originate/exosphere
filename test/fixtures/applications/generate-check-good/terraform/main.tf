variable "aws_profile" {
  default = "default"
}

terraform {
  required_version = ">= 0.10.0"

  backend "s3" {
    bucket         = "-out-of-date-yaml-terraform"
    key            = "terraform.tfstate"
    region         = ""
    dynamodb_table = "TerraformLocks"
  }
}

provider "aws" {
  version = "0.1.4"

  region              = ""
  profile             = "${var.aws_profile}"
  allowed_account_ids = [""]
}

variable "key_name" {
  default = ""
}

module "aws" {
  source = "git@github.com:Originate/exosphere.git//terraform//aws?ref=1bf0375f"

  name              = "out-of-date-yaml"
  env               = "production"
  external_dns_name = ""
  key_name          = "${var.key_name}"
}

variable "test-service_env_vars" {
  default = "[]"
}

variable "test-service_docker_image" {}

module "test-service" {
  source = "git@github.com:Originate/exosphere.git//terraform//aws//worker-service?ref=1bf0375f"

  name = "test-service"

  cluster_id            = "${module.aws.ecs_cluster_id}"
  cpu                   = ""
  desired_count         = 1
  docker_image          = "${var.test-service_docker_image}"
  env                   = "production"
  environment_variables = "${var.test-service_env_vars}"
  memory_reservation    = ""
  region                = "${module.aws.region}"
}

