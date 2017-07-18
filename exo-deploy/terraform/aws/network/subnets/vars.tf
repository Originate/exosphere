/* Variables */

variable "availability_zones" {
  description = "Availability zones to use for subnets. Two subnets will be created per availability zone"
  type        = "list"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "name" {
  description = "Name tag, e.g stack"
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

/* Output */

output "private_subnet_ids" {
  description = "ID's of the private subnets"
  value       = ["${aws_subnet.private.*.id}"]
}

output "public_subnet_ids" {
  description = "ID's of the public subnets"
  value       = ["${aws_subnet.public.*.id}"]
}
