/* Variables */

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "logs_expiration_days" {
  default = 30
}

variable "logs_expiration_enabled" {
  default = false
}

variable "bucket_prefix" {
  description = "Bucket prefix, e.g stack"
}

/* Output */

output "id" {
  description = "S3 bucket ID"
  value       = "${aws_s3_bucket.logs.id}"
}
