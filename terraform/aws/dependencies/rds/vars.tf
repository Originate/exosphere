/* Variables */

variable "allocated_storage" {
  description = "Allocated storage in gigabytes"
}


variable "ecs_security_group" {
  description = "Security group ID of ECS cluster from which to allow traffic"
}

variable "bastion_security_group" {
  description = "Security group ID of bastion instances to allow ssh connection"
}

variable "engine" {
  description = "Database engine to use"
}

variable "engine_version" {
  description = "Database engine version"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "identifier" {
  description = "Name of RDS instance"
  default     = ""
}

variable "instance_class" {
  description = "Instance type of RDS instance"
}

variable "internal_hosted_zone_id" {
  description = "Route53 Hosted Zone id used for internal routing"
}

variable "name" {
  description = "Name of database to create when RDS instance created"
}

variable "username" {
  description = "Username for master db user."
}

variable "password" {
  description = "Password for master db user. Note that this may show up in logs, and it will be stored in the state file."
}

variable "storage_type" {
  description = "Storage type, i.e. general purpose SSD, provisioned IOPS, magnetic, etc."
}

variable "subnet_ids" {
  description = "List of subnet IDs the ALB should live in"
  type        = "list"
}

variable "vpc_id" {
  description = "ID of the main VPC"
}

/* Output */

output "endpoint" {
  description = "Connection endpoint"
  value       = "${aws_db_instance.rds.endpoint}"
}
