variable "availability_zones" {
  description = "List of AZs"
  type        = "list"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "key_name" {
  description = "Name of the key pair stored in AWS used to SSH into bastion instances"
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
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

output "bastion_security_group_id" {
  value = "${module.bastion.security_group_id}"
}
