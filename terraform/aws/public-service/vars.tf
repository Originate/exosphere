/* Variables */

variable "alb_security_group" {
  description = "ID of external ALB security group"
}

variable "alb_subnet_ids" {
  description = "List of public subnet ID's the ALB should live in"
  type        = "list"
}

variable "cluster_id" {
  description = "ID of the ECS cluster"
}

variable "command" {
  description = "Starting command to run in container"
  type        = "list"
  default     = []
}

variable "container_port" {
  description = "Port number on the container to bind the ALB to"
}

variable "cpu" {
  description = "Number of cpu units to reserve for the container"
}

variable "docker_image" {
  description = "ECS repository URI of Docker image"
}

variable "desired_count" {
  description = "Desired number of tasks to keep running"
  default = 2
}

variable "ecs_role_arn" {
  description = "ARN of the ECS IAM role"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "environment_variables" {
  description = "Environment variables to pass to a container"
  default     = "[]"
}

variable "external_dns_name" {
  description = "External DNS name to host public service at"
}

variable "external_zone_id" {
  description = "Route53 Hosted Zone id used for external routing"
}

variable "health_check_endpoint" {
  description = "Endpoint for the alb to hit when performing health checks"
  default     = "/"
}

variable "internal_dns_name" {
  description = "Internal DNS name used for internal routing"
}

variable "internal_zone_id" {
  description = "Route53 Hosted Zone id used for internal routing"
}

variable "log_bucket" {
  description = "S3 bucket id to write ELB logs into"
}

variable "memory_reservation" {
  description = "Soft limit (in MiB) of memory to reserve for the container"
}

variable "name" {
  description = "Name of the service"
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

variable "ssl_certificate_arn" {
  description = "The ARN of the SSL server certificate"
}

variable "vpc_id" {
  description = "ID of the VPC"
}
