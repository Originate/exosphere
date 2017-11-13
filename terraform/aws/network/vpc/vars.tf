/* Variables */

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "name" {
  description = "Name tag, e.g stack"
}

variable "cidr" {
  description = "The CIDR block of the VPC"
  default     = "10.0.0.0/16"
}

/* Output */

output "id" {
  description = "ID of the VPC"
  value       = "${aws_vpc.vpc.id}"
}

output "cidr" {
  description = "The CIDR block of the VPC"
  value       = "${aws_vpc.vpc.cidr_block}"
}
