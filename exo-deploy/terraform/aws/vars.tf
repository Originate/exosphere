/* Variables */

variable "ecs_ebs_optimized" {
  description = "Boolean indicating if ECS instances are EBS-optimized"
  default     = false
}

variable "ecs_instance_type" {
  description = "Instance type to use for ECS instances"
  default     = "t2.micro"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "key_name" {
  description = "Name of the key pair stored in AWS used to SSH into bastion instances"
}

variable "name" {
  description = "Application name"
}

/* Outputs */

output "availability_zones" {
  description = "List of AZs"
  value       = ["${data.aws_availability_zones.available.names}"]
}

output "bastion_ips" {
  description = "IP addresses of the bastion hosts"
  value       = ["${module.network.bastion_ips}"]
}

output "bastion_security_group" {
  description = "ID of the security group of the bastion hosts"
  value       = "${module.network.bastion_security_group}"
}

output "ecs_cluster_id" {
  description = "ID of the ECS cluster"
  value       = "${module.ecs_cluster.id}"
}

output "ecs_cluster_security_group" {
  description = "ID of the security group of the ECS cluster instances"
  value       = "${module.ecs_cluster.security_group}"
}

output "ecs_service_iam_role_arn" {
  description = "ARN of ECS service IAM role passed to each service module"
  value       = "${module.ecs_cluster.ecs_service_iam_role_arn}"
}

output "external_alb_security_group" {
  description = "ID of the external ALB security group"
  value       = "${module.alb_security_groups.external_alb_id}"
}

output "internal_alb_security_group" {
  description = "ID of the internal ALB security group"
  value       = "${module.alb_security_groups.internal_alb_id}"
}

output "internal_zone_id" {
  description = "The Route53 internal zone ID"
  value       = "${module.internal_dns.zone_id}"
}

output "log_bucket_id" {
  description = "S3 bucket id of load balancer logs"
  value       = "${module.s3_logs.id}"
}

output "public_subnet_ids" {
  description = "ID's of the public subnets"
  value       = ["${module.network.public_subnet_ids}"]
}

output "private_subnet_ids" {
  description = "ID's of the private subnets"
  value       = ["${module.network.private_subnet_ids}"]
}

output "region" {
  description = "Region of the environment, for example, us-west-2"
  value       = "${data.aws_region.current.name}"
}

output "vpc_id" {
  value = "${module.network.vpc_id}"
}
