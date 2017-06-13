variable "env" {
  description = "Environment tag, e.g prod"
}

variable "availability_zones" {
  description = "List of AZs"
  type        = "list"
}

variable "region" {}

variable "key_name" {
  description = "Name of the SSH key pair stored in AWS to authorize for the bastion hosts"
}
