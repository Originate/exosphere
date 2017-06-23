/* Variables */

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "health_check_endpoint" {
  description = "Endpoint for the alb to hit when performing health checks"
  default     = "/"
}

variable "name" {
  description = "Name of the service"
}

variable "security_groups" {
  description = "ID of external ALB security group"
  type        = "list"
}

variable "subnet_ids" {
  description = "List of public or private ID's the ALB should live in"
  type        = "list"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

/* Output */

output "url" {
  value = "${aws_alb.alb.dns_name}"
}

output "target_group_id" {
  value = "${aws_alb_target_group.target_group.id}"
}
