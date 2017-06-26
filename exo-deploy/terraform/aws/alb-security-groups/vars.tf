/* Variables */

variable "name" {
  description = "Name of application"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "vpc_id" {
  description = "VPC ID"
}

/* Output */

output "internal_alb_security_group" {
  description = "ID of internal ALB security groups"
  value       = "${aws_security_group.internal_alb.id}"
}

output "external_alb_security_group" {
  description = "ID of external ALB security groups"
  value       = "${aws_security_group.external_alb.id}"
}
