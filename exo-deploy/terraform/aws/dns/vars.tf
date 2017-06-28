/* Variables */

variable "name" {
  description = "Zone name, e.g stack.local"
}

variable "vpc_id" {
  description = "The VPC ID (omit to create a public zone)"
  default     = ""
}

/* Output */

output "name" {
  description = "The domain name"
  value       = "${var.name}"
}

output "zone_id" {
  description = "The zone ID"
  value = "${aws_route53_zone.main.zone_id}"
}

output "name_servers" {
  description = "A comma separated list of the zone name servers"
  value = "${join(",",aws_route53_zone.main.name_servers)}"
}
