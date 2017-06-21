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

variable "ebs_optimized" {
  description = "Boolean indicating if cluster instances are ebs optimized"
  default     = "false"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "iam_instance_profile" {
  description = "IAM instance profile passed to the cluster launch configuration"
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

variable "high_cpu_threshold" {
  description = "If CPU usage is above this threshold for 5min, scale up"
  default     = 90
}

variable "low_cpu_threshold" {
  description = "If CPU usage is below this threshold for 5min, scale down"
  default     = 10
}

variable "high_memory_threshold" {
  description = "If CPU usage is above this threshold for 5min, scale up"
  default     = 90
}

variable "low_memory_threshold" {
  description = "If CPU usage is below this threshold for 5min, scale down"
  default     = 10
}

variable "docker_auth_type" {
  description = "The docker auth type, see https://godoc.org/github.com/aws/amazon-ecs-agent/agent/engine/dockerauth for the possible values"
  default     = ""
}

variable "docker_auth_data" {
  description = "A JSON object providing the docker auth data, see https://godoc.org/github.com/aws/amazon-ecs-agent/agent/engine/dockerauth for the supported formats"
  default     = ""
}

variable "extra_cloud_config_type" {
  description = "Extra cloud config type"
  default     = "text/cloud-config"
}

variable "extra_cloud_config_content" {
  description = "Extra cloud config content"
  default     = ""
}
