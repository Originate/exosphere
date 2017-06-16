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

variable "subnet_ids" {
  type        = "list"
  description = "List of public or private ID's the ALB should live in"
}

variable "vpc_id" {
  description = "ID of the VPC"
}
