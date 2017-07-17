/* Variables */

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "name" {
  description = "Zone name, e.g stack.local"
}

variable "servers" {
  description = "List of the IP addresses of internal DNS servers"
  type        = "list"
}

variable "vpc_id" {
  description = "The VPC ID"
}

/* Output */

output "zone_id" {
  description = "The zone ID"
  value       = "${aws_route53_zone.zone.zone_id}"
}
