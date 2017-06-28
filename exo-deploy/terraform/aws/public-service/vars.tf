/* Variables */

variable "alb_security_group" {
  description = "ID of external ALB security group"
  type        = "list"
}

variable "alb_subnet_ids" {
  type        = "list"
  description = "List of public or private ID's the ALB should live in"
}

variable "cluster_id" {
  description = "ID of the ECS cluster"
}

variable "command" {
  description = "Starting command to run in container"
  type = "list"
}

variable "container_port" {
  description = "Port number on the container to bind the ALB to"
  default     = 80
}

variable "cpu_units" {
  description = "Number of cpu units to reserve for the container"
}

variable "docker_image" {
  description = "ECS repository URI of Docker image"
}

variable "domain_name" {
  description = "Domain name to host ExoCom at"
}

variable "ecs_role_arn" {
  description = "ARN of the ECS IAM role"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "environment_variables" {
  type        = "map"
  description = "Environment variables to pass to a container"
}

variable "health_check_endpoint" {
  description = "Endpoint for the alb to hit when performing health checks"
  default     = "/"
}

variable "hosted_zone_id" {
  description = "Route53 Hosted Zone id with registered NS records"
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

variable "vpc_id" {
  description = "ID of the VPC, passed to ALB"
}


/* Output */

output "url" {
  value = "${aws_route53_record.public_url.name}"
}
