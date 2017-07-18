/* Variables */

variable "name" {
  description = "Name of application"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

/* Output */

output "internal_alb_id" {
  description = "ID of internal ALB security group"
  value       = "${aws_security_group.internal_alb.id}"
}

output "external_alb_id" {
  description = "ID of external ALB security group"
  value       = "${aws_security_group.external_alb.id}"
}
