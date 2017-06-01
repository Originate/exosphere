variable "env" {
  description = "Environment tag, e.g prod"
}

variable "availability_zones" {
  description = "List of AZs"
  type        = "list"
}

module "vpc" {
  source = "./vpc"

  env = "${var.env}"
}

module "subnets" {
  source = "./subnets"

  env                = "${var.env}"
  vpc_id             = "${module.vpc.id}"
  vpc_cidr           = "${module.vpc.cidr}"
  availability_zones = ["${var.availability_zones}"]
}

output "vpc_id" {
  value = "${module.vpc.id}"
}

output "public_subnet_ids" {
  value = ["${module.subnets.public_subnet_ids}"]
}

output "private_subnet_ids" {
  value = ["${module.subnets.private_subnet_ids}"]
}
