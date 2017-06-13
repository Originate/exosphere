variable "env" {}

variable "region" {}

variable "account_id" {
  default = ""
}

variable "security_groups" {
  description = "Comma separated list of security groups"
  type        = "list"
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

module "network" {
  source = "./network"

  env                = "${var.env}"
  availability_zones = "${data.aws_availability_zones.available.names}"
  region             = "${var.region}"
  key_name           = "hugo"
}

module "cluster" {
  source = "./cluster"

  name       = "exocom-cluster"
  env        = "${var.env}"
  vpc_id     = "${module.network.vpc_id}"
  subnet_ids = ["${module.network.private_subnet_ids}"]

  security_groups      = ["${var.security_groups}","${module.network.bastion_security_group_id}"]
  /* iam_instance_profile = "${module.iam.iam_instance_profile}" */
  availability_zones   = "${data.aws_availability_zones.available.names}"

  image_id      = "${data.aws_ami.ecs_optimized_ami.id}"
  instance_type = "t2.micro"
  key_name = "hugo"
}

output "vpc_id" {
  value = "${module.network.vpc_id}"
}

output "public_subnet_ids" {
  value = ["${module.network.public_subnet_ids}"]
}

output "cluster_id" {
  value = "${module.cluster.id}"
}

output "cluster_security_group_id" {
  value = "${module.cluster.security_group_id}"
}
