/* Variables */

variable "availability_zones" {
  type        = "list"
  description = "List of availability zones to place subnets"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "instance_type" {
  description = "Instance type of the bastion hosts"
  default     = "t2.micro"
}

variable "key_name" {
  description = "Name of the key pair stored in AWS used to SSH into bastion instances"
}

variable "name" {
  description = "Name tag, e.g stack"
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

variable "subnet_ids" {
  description = "ID's of the subnets to place instances into"
  type        = "list"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

/* Output */

output "ips" {
  description = "IP addresses of the bastion hosts"
  value       = ["${formatlist("ubuntu@%s", aws_instance.bastion.*.public_ip)}"]
}

output "security_group" {
  description = "ID of the security group of the bastion hosts"
  value       = "${aws_security_group.bastion.id}"
}
