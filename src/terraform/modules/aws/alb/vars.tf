/* Variables */

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "external_dns_name" {
  description = "The subdomain under which the ALB is exposed externally. This is not used for internal ALBs"
  default     = ""
}

variable "external_zone_id" {
  description = "The Route53 zone ID to create the external record in. This is not used for internal ALBs"
  default     = ""
}

variable "health_check_endpoint" {
  description = "Endpoint for the alb to hit when performing health checks"
  default     = "/"
}

variable "internal" {
  description = "Set this to false for public ALBs"
  default     = true
}

variable "internal_dns_name" {
  description = "Internal DNS name used for routing"
}

variable "internal_zone_id" {
  description = "Hosted zone ID used for internal routing"
}

variable "log_bucket" {
  description = "S3 bucket id to write ELB logs into"
}

variable "name" {
  description = "Name tag, e.g stack"
}

variable "security_groups" {
  description = "IDs of the ALBs security groups"
  type        = "list"
}

variable "ssl_certificate_arn" {
  description = "The ARN of the SSL server certificate. This is not used for internal ALBs"
  default     = ""
}

variable "subnet_ids" {
  description = "List of subnet IDs the ALB should live in"
  type        = "list"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

/* Output */

output "dns_name" {
  description = "The DNS name"
  value       = "${aws_alb.alb.dns_name}"
}

output "zone_id" {
  description = "The canonical hosted zone ID"
  value       = "${aws_alb.alb.zone_id}"
}

output "target_group_id" {
  description = "ID of the target group"
  value       = "${aws_alb_target_group.target_group.id}"
}
