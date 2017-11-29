variable "aws_profile" {
  default = "default"
}

terraform {
  required_version = "= 0.11.0"

  backend "s3" {
    bucket         = "-clean-containers-terraform"
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

  name              = "clean-containers"
  env               = "production"
  external_dns_name = ""
  key_name          = "${var.key_name}"
}

variable "application-service_env_vars" {
  default = "[]"
}

variable "application-service_docker_image" {}

module "application-service" {
  source = "git@github.com:Originate/exosphere.git//terraform//aws//public-service?ref=1bf0375f"

  name = "application-service"

  alb_security_group    = "${module.aws.external_alb_security_group}"
  alb_subnet_ids        = ["${module.aws.public_subnet_ids}"]
  cluster_id            = "${module.aws.ecs_cluster_id}"
  container_port        = ""
  cpu                   = ""
  desired_count         = 1
  docker_image          = "${var.application-service_docker_image}"
  ecs_role_arn          = "${module.aws.ecs_service_iam_role_arn}"
  env                   = "production"
  environment_variables = "${var.application-service_env_vars}"
  external_dns_name     = ""
  external_zone_id      = "${module.aws.external_zone_id}"
  health_check_endpoint = ""
  internal_dns_name     = "application-service"
  internal_zone_id      = "${module.aws.internal_zone_id}"
  log_bucket            = "${module.aws.log_bucket_id}"
  memory_reservation    = ""
  region                = "${module.aws.region}"
  ssl_certificate_arn   = ""
  vpc_id                = "${module.aws.vpc_id}"
}

