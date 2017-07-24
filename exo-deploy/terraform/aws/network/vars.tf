/* Variables */

variable "name" {
  description = "Name tag, e.g stack"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "availability_zones" {
  description = "List of AZs"
  type        = "list"
}

variable "key_name" {
  description = "Name of the key pair stored in AWS used to SSH into bastion instances"
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

/* Output */

output "vpc_cidr" {
  description = "The CIDR block of the VPC"
  value       = "${module.vpc.cidr}"
}

output "vpc_id" {
  description = "ID of the VPC"
  value       = "${module.vpc.id}"
}

output "public_subnet_ids" {
  description = "ID's of the public subnets"
  value       = ["${module.subnets.public_subnet_ids}"]
}

output "private_subnet_ids" {
  description = "ID's of the private subnets"
  value       = ["${module.subnets.private_subnet_ids}"]
}

output "ssh_config" {
  description = "SSH config of the bastion hosts"
  value       = ["${module.bastion.ips}"]
}

output "bastion_security_group" {
  description = "ID of the security group of the bastion hosts"
  value       = "${module.bastion.security_group}"
}
