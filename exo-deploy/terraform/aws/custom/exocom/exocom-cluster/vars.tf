/* Variable */

variable "availability_zones" {
  description = "List of AZs"
  type        = "list"
}

variable "desired_capacity" {
  description = "Desired instance count"
  default     = 3
}

variable "docker_auth_data" {
  description = "A JSON object providing the docker auth data, see https://godoc.org/github.com/aws/amazon-ecs-agent/agent/engine/dockerauth for the supported formats"
  default     = ""
}

variable "docker_auth_type" {
  description = "The docker auth type, see https://godoc.org/github.com/aws/amazon-ecs-agent/agent/engine/dockerauth for the possible values"
  default     = ""
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

variable "extra_cloud_config_content" {
  description = "Extra cloud config content"
  default     = ""
}

variable "extra_cloud_config_type" {
  description = "Extra cloud config type"
  default     = "text/cloud-config"
}

variable "high_cpu_threshold" {
  description = "If CPU usage is above this threshold for 5min, scale up"
  default     = 90
}

variable "high_memory_threshold" {
  description = "If CPU usage is above this threshold for 5min, scale up"
  default     = 90
}

variable "instance_type" {
  description = "The instance type to use, e.g t2.small"
}

variable "key_name" {
  description = "Name of key pair stored in AWS to authorize for the bastion hosts"
}

variable "low_cpu_threshold" {
  description = "If CPU usage is below this threshold for 5min, scale down"
  default     = 10
}

variable "low_memory_threshold" {
  description = "If CPU usage is below this threshold for 5min, scale down"
  default     = 10
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

variable "region" {
  description = "Region of the environment, for example, us-west-2"
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

/* Output */

output "cluster_id" {
  description = "ID of main cluster"
  value       = "${aws_ecs_cluster.exocom.id}"
}

output "security_groups" {
  description = "Cluster and external alb sg ids"
  value       = ["${aws_security_group.exocom_cluster.id}", "${aws_security_group.external_alb.id}"]
}
