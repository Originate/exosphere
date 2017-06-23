/* Variables */

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "health_check_endpoint" {
  description = "Endpoint for the alb to hit when performing health checks"
  default     = "/"
}

variable "host_port" {
  description = "Port number on the host to bind the container to"
  default     = 80
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

output "name" {
  value = "${aws_elb.elb.name}"
}

output "url" {
  value = "${aws_elb.elb.dns_name}"
}
