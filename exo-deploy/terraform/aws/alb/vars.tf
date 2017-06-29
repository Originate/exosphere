/* Variables */

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "health_check_endpoint" {
  description = "Endpoint for the alb to hit when performing health checks"
  default     = "/"
}

variable "log_bucket" {
  description = "S3 bucket id to write ELB logs into"
}

variable "name" {
  description = "Name of the service"
}

variable "security_group" {
  description = "ID of ALB security group"
  type        = "list"
}

variable "subnet_ids" {
  type        = "list"
  description = "List of public or private ID's the ALB should live in"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

/* Output */

output "dns_name" {
  value = "${aws_alb.alb.dns_name}"
}

output "zone_id"  {
  value = "${aws_alb.alb.zone_id}"
}

output "target_group_id" {
  value = "${aws_alb_target_group.target_group.id}"
}
