variable "aws_profile" {
  default = "default"
}

variable "aws_region" {}

variable "aws_account_id" {}

variable "aws_ssl_certificate_arn" {}

variable "application_url" {}

variable "env" {}

terraform {
  required_version = "= {{{terraformVersion}}}"

  backend "s3" {
    key            = "infrastructure.tfstate"
    dynamodb_table = "{{lockTable}}"
  }
}

provider "aws" {
  version = "0.1.4"

  region              = "${var.aws_region}"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["${var.aws_account_id}"]
}

variable "key_name" {
  default = ""
}

module "aws" {
  source = "github.com/Originate/exosphere.git//terraform//aws?ref={{terraformCommitHash}}"

  name              = "{{appName}}"
  env               = "${var.env}"
  external_dns_name = "${var.application_url}"
  key_name          = "${var.key_name}"
  log_bucket_prefix = "${var.aws_account_id}-{{appName}}-${var.env}"
}

output "availability_zones" {
  description = "List of AZs"
  value       = "${module.aws.availability_zones}"
}

output "bastion_security_group" {
  description = "ID of the security group of the bastion hosts"
  value       = "${module.aws.bastion_security_group}"
}

output "ecs_cluster_id" {
  description = "ID of the ECS cluster"
  value       = "${module.aws.ecs_cluster_id}"
}

output "ecs_cluster_security_group" {
  description = "ID of the security group of the ECS cluster instances"
  value       = "${module.aws.ecs_cluster_security_group}"
}

output "ecs_service_iam_role_arn" {
  description = "ARN of ECS service IAM role passed to each service module"
  value       = "${module.aws.ecs_service_iam_role_arn}"
}

output "external_alb_security_group" {
  description = "ID of the external ALB security group"
  value       = "${module.aws.external_alb_security_group}"
}

output "external_zone_id" {
  description = "The Route53 external zone ID"
  value       = "${module.aws.external_zone_id}"
}

output "internal_alb_security_group" {
  description = "ID of the internal ALB security group"
  value       = "${module.aws.internal_alb_security_group}"
}

output "internal_zone_id" {
  description = "The Route53 internal zone ID"
  value       = "${module.aws.internal_zone_id}"
}

output "log_bucket_id" {
  description = "S3 bucket id of load balancer logs"
  value       = "${module.aws.log_bucket_id}"
}

output "public_subnet_ids" {
  description = "ID's of the public subnets"
  value       = ["${module.aws.public_subnet_ids}"]
}

output "private_subnet_ids" {
  description = "ID's of the private subnets"
  value       = ["${module.aws.private_subnet_ids}"]
}

output "region" {
  description = "Region of the environment, for example, us-west-2"
  value       = "${module.aws.region}"
}

output "vpc_id" {
  value = "${module.aws.vpc_id}"
}
