variable "cluster_id" {
  description = "ID of the ECS cluster"
}

variable "command" {
  description = "Starting command to run in container"
  type        = "list"
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

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "environment_variables" {
  type        = "map"
  description = "Environment variables to pass to a container"
}

variable "memory" {
  description = "Soft limit (in MiB) of memory to reserve for the container"
}

variable "name" {
  description = "Name of the service"
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
}
