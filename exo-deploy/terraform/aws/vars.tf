/* Variables */

variable "account_id" {
  description = "ID associated with AWS account"
  default     = ""
}

variable "application_name" {
  description = "Application name (kabob case)"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "key_name" {
  description = "Name of the key pair stored in AWS used to SSH into bastion instances"
  default     = ""
}

variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

variable "security_groups" {
  description = "Comma separated list of security groups passed to main cluster"
  type        = "list"
}


/* Outputs */

output "availability_zones" {
  value       = "${data.aws_availability_zones.available.names}"
  description = "Names of availability zones"
}

output "bastion_security_group_id" {
  value       = "${module.network.bastion_security_group_id}"
  description = "Security group ID of the bastion instances used to ssh into cluster instances"
}

output "cluster_id" {
  value = "${module.cluster.id}"
}

output "cluster_security_group" {
  value       = "${module.cluster.security_group}"
  description = "ID of security group for main cluster, passed to each service."
}

output "ecs_service_iam_role_arn" {
  value       = "${module.cluster.ecs_service_iam_role_arn}"
  description = "ARN of ECS service IAM role passed to each service module"
}

output "external_alb_security_group" {
  description = "ID of external ALB security groups"
  value       = "${module.alb_security_groups.external_alb_security_group}"
}

output "internal_alb_security_group" {
  description = "ID of internal ALB security groups"
  value       = "${module.alb_security_groups.internal_alb_security_group}"
}

output "internal_dns_name" {
  description = "Internal DNS name used for internal routing"
  value       = "${module.dns.name}"
}

output "internal_hosted_zone_id" {
  description = "Hosted zone id created by used for internal routing"
  value       = "${module.dns.zone_id}"
}

output "log_bucket_id" {
  description = "S3 bucket id of load balancer logs"
  value       = "${module.s3_logs.id}"
}

output "public_subnet_ids" {
  value = ["${module.network.public_subnet_ids}"]
}

output "private_subnet_ids" {
  value = ["${module.network.private_subnet_ids}"]
}

output "vpc_id" {
  value = "${module.network.vpc_id}"
}
