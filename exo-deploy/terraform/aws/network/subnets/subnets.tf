variable "env" {
  description = "Environment tag, e.g prod"
}

variable "vpc_id" {}

variable "vpc_cidr" {
  description = "The cidr block of the VPC"
}

variable "availability_zones" {
  description = "Availability zones to use for subnets. Two subnets will be created per availability zone"
  type        = "list"
}
