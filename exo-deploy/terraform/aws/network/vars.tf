variable "availability_zones" {
  description = "List of AZs"
  type        = "list"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "key_name" {
  description = "Name of the key pair stored in AWS used to SSH into bastion instances"
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

