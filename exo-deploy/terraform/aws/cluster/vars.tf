variable "availability_zones" {
  description = "List of AZs"
  type        = "list"
}

variable "desired_capacity" {
  description = "Desired instance count"
  default     = 3
}

variable "docker_volume_size" {
  description = "Attached EBS volume size in GB"
  default     = 25
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "iam_instance_profile" {
  description = "IAM instance profile passed to the cluster launch configuration"
}

variable "image_id" {
  description = "AMI Image ID"
}

variable "instance_type" {
  description = "The instance type to use, e.g t2.small"
}

variable "key_name" {
  description = "Name of key pair stored in AWS to authorize for the bastion hosts"
}

variable "max_size" {
  description = "Maxmimum instance count"
  default     = 100
}

variable "min_size" {
  description = "Minimum instance count"
  default     = 3
}

variable "name" {
  description = "The cluster name, e.g cdn"
}

variable "root_volume_size" {
  description = "Root volume size in GB"
  default     = 25
}

variable "security_groups" {
  description = "Comma separated list of security groups"
  type        = "list"
}

variable "subnet_ids" {
  description = "List of subnet IDs"
  type        = "list"
}

variable "vpc_id" {
  description = "VPC ID"
}
