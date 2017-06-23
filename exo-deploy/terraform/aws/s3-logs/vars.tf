/* Variables */

variable "name" {
  description = "Name of application"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "account_id" {
  description = "ID associated with AWS account"
  default     = ""
}

variable "logs_expiration_enabled" {
  default = false
}

variable "logs_expiration_days" {
  default = 30
}

/* Output */

output "id" {
  description = "S3 bucket ID"
  value       = "${aws_s3_bucket.logs.id}"
}
