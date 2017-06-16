variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

variable "account_id" {
  description = "ID associated with AWS account"
  default     = ""
}

variable "security_groups" {
  description = "Comma separated list of security groups passed to main cluster"
  type        = "list"
}

variable "key_name" {
  description = "Name of the key pair stored in AWS used to SSH into bastion instances"
  default     = ""
}

data "aws_availability_zones" "available" {}

/* Get ECS optimized AMI id to use on the cluster */
data "aws_ami" "ecs_optimized_ami" {
  most_recent = true

  filter {
    name   = "owner-alias"
    values = ["amazon"]
  }

  filter {
    name   = "name"
    values = ["amzn-ami-2016.09.g-amazon-ecs-optimized"]
  }
}

provider "aws" {
  region              = "${var.region}"
  allowed_account_ids = ["${var.account_id}"]
}

module "iam" {
  source = "./iam"
  env    = "${var.env}"
}

module "network" {
  source             = "./network"

  env                = "${var.env}"
  availability_zones = "${data.aws_availability_zones.available.names}"
  region             = "${var.region}"
  key_name           = "${var.key_name}"
}

module "cluster" {
  source               = "./cluster"

  name                 = "exosphere-cluster"
  env                  = "${var.env}"
  vpc_id               = "${module.network.vpc_id}"
  subnet_ids           = ["${module.network.private_subnet_ids}"]

  iam_instance_profile = "${module.iam.iam_instance_profile}"
  security_groups      = ["${module.network.bastion_security_group_id}", "${var.security_groups}"]
  availability_zones   = "${data.aws_availability_zones.available.names}"

  image_id             = "${data.aws_ami.ecs_optimized_ami.id}"
  instance_type        = "t2.micro"
  key_name             = "${var.key_name}"
}

output "vpc_id" {
  value = "${module.network.vpc_id}"
}

output "public_subnet_ids" {
  value = ["${module.network.public_subnet_ids}"]
}

output "private_subnet_ids" {
  value = ["${module.network.private_subnet_ids}"]
}

output "cluster_id" {
  value = "${module.cluster.id}"
}

output "cluster_security_group_id" {
  value       = "${module.cluster.security_group_id}"
  description = "ID of main cluster passed to each service module"
}

output "ecs_iam_role_arn" {
  value       = "${module.iam.iam_role_arn}"
  description = "ARN of ECS IAM role passed to each service module"
}
