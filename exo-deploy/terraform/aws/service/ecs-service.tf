variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "name" {
  description = "Name of the service"
}

variable "command" {
  description = "Command to run in container"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

variable "public_subnet_ids" {
  type        = "list"
  description = "List of ID's of the public subnets"
}

/* variable "cluster_security_group_id" { */
/*   type        = "list" */
/*   description = "Cluster security group id" */
/* } */

variable "cluster_id" {
  description = "ID of the ECS cluster"
}

/* variable "ecs_role_arn" { */
/*   description = "ARN of the ECS IAM role" */
/* } */

variable "cpu_units" {
  description = "Number of cpu units to reserve for the container"
}

variable "memory_reservation" {
  description = "Soft limit (in MiB) of memory to reserve for the container"
}

variable "docker_image" {
  description = "Docker image to use"
}

variable "container_port" {
  description = "Port number on the container to bind"
}

variable "environment_variables" {
  type        = "map"
  description = "Environment variables to pass to a container"
}
