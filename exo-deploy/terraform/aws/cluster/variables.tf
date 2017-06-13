variable "name" {
  description = "The cluster name, e.g cdn"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "vpc_id" {
  description = "VPC ID"
}

variable "subnet_ids" {
  description = "List of subnet IDs"
  type        = "list"
}

variable "security_groups" {
  description = "Comma separated list of security groups"
  type        = "list"
}

variable "key_name" {
  description = "Name of the SSH key pair stored in AWS to authorize for the bastion hosts"
}

variable "availability_zones" {
  description = "List of AZs"
  type        = "list"
}

variable "image_id" {
  description = "AMI Image ID"
}

variable "instance_type" {
  description = "The instance type to use, e.g t2.small"
}

variable "min_size" {
  description = "Minimum instance count"
  default     = 3
}

variable "max_size" {
  description = "Maxmimum instance count"
  default     = 100
}

variable "desired_capacity" {
  description = "Desired instance count"
  default     = 3
}

variable "root_volume_size" {
  description = "Root volume size in GB"
  default     = 25
}

variable "docker_volume_size" {
  description = "Attached EBS volume size in GB"
  default     = 25
}

