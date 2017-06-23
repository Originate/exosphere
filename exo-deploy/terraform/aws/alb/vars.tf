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

variable "security_group" {
  description = "ID of ALB security group"
}

variable "subnet_ids" {
  type        = "list"
  description = "List of public or private ID's the ALB should live in"
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
