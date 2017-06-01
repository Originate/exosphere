variable "env" {}

variable "region" {}

variable "account_profile" {}

variable "account_id" {
  default = ""
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
  profile             = "${var.account_profile}"
  allowed_account_ids = ["${var.account_id}"]
}

module "network" {
  source = "./network"

  env                = "${var.env}"
  availability_zones = ["${data.aws_availability_zones.available.names}"]
}

module "cluster" {
  source = "./cluster"

  name       = "exocom-cluster"
  env        = "${var.env}"
  vpc_id     = "${module.network.vpc_id}"
  subnet_ids = ["${module.network.private_subnet_ids}"]

  security_groups      = []                                                 // TODO
  iam_instance_profile = ""                                                 // TODO
  availability_zones   = ["${data.aws_availability_zones.available.names}"]

  image_id      = "${data.aws_ami.ecs_optimized_ami.id}"
  instance_type = "t2.micro"
}
