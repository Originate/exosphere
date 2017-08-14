/* Variables */

variable "name" {
  description = "The cluster name"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

variable "alb_security_groups" {
  description = "List of ID's of the security groups of the ALB's"
  type        = "list"
}

variable "bastion_security_group" {
  description = "ID of the security group of the bastion hosts"
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
  description = "Boolean indicating if cluster instances are EBS-optimized"
  default     = false
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
  description = "If memory usage is above this threshold for 5min, scale up"
  default     = 90
}

variable "low_cpu_threshold" {
  description = "If CPU usage is below this threshold for 5min, scale down"
  default     = 10
}

variable "low_memory_threshold" {
  description = "If memory usage is below this threshold for 5min, scale down"
  default     = 10
}

variable "desired_capacity" {
  description = "Desired instance count"
  default     = 3
}

variable "max_size" {
  description = "Maxmimum instance count"
  default     = 100
}

variable "min_size" {
  description = "Minimum instance count"
  default     = 3
}

variable "instance_type" {
  description = "The instance type to use, e.g t2.small"
}

variable "key_name" {
  description = "Name of key pair stored in AWS to authorize for the bastion hosts"
}


variable "root_volume_size" {
  description = "Root volume size in GB"
  default     = 25
}

variable "subnet_ids" {
  description = "List of subnet IDs that cluster lives in"
  type        = "list"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

/* Output */

output "id" {
  description = "ID of cluster"
  value       = "${aws_ecs_cluster.cluster.id}"
}

output "ecs_service_iam_role_arn" {
  description = "ARN of ECS service IAM role"
  value       = "${aws_iam_role.ecs_service.arn}"
}

output "security_group" {
  description = "Cluster security group ID"
  value       = "${aws_security_group.cluster.id}"
}
