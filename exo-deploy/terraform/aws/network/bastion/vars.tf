/* Variables */

variable "availability_zones" {
  type        = "list"
  description = "List of availability zones to place subnets"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "instance_type" {
  default     = "t2.micro"
  description = "Instance type of the bastion hosts"
}

variable "key_name" {
  description = "Name of the key pair stored in AWS used to SSH into bastion instances"
}

variable "public_subnet_ids" {
  type        = "list"
  description = "List of ID's of the public subnets"
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

/* Output */

output "private_ips" {
  value = ["${aws_instance.bastion.*.private_ip}"]
}

output "public_ips" {
  value = ["${aws_instance.bastion.*.public_ip}"]
}

output "security_group_id" {
  value = "${aws_security_group.bastion.id}"
}

