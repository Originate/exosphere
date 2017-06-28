/* Variables */

variable "command" {
  description = "Starting command to run in container"
  type = "list"
}

variable "container_port" {
  description = "Port number on the container to bind the host to"
  default     = 80
}

variable "cpu_units" {
  description = "Number of cpu units to reserve for the container"
}

variable "docker_image" {
  description = "ECS repository URI of Docker image"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "environment_variables" {
  type        = "map"
  description = "Environment variables to pass to a container"
}

variable "host_port" {
  description = "Port number on the host to bind the container to"
  default     = 80
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

/* Output */

output "task_arn" {
  value       = "${aws_ecs_task_definition.task.arn}"
  description = "ARN of task definition to be passed to ECS service"
}
