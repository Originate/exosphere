variable "env" {
  description = "Environment tag, e.g prod"
}

variable "cidr" {
  description = "The cidr block for the VPC"
  default     = "10.0.0.0/16"
}

variable "name" {
  default = "vpc"
}
