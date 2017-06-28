variable "access_key" {
  description = "AWS access key"
}

variable "account_id" {
  description = "ID of AWS account"
}

variable "hosted_zone_id" {
  description = "Route53 Hosted Zone id with registered NS records"
}

variable "key_name" {
  description = "Key pair name used for SSH"
}

variable "mongodb_pw" {
  description = "Environment variable for mlabs password. Prompted for during 'terraform plan/apply' or set using TF_VAR_MONGODB_PW=#{value}"
}

variable "mongodb_user" {
  description = "Environment variable for mlabs username. Prompted for during 'terraform plan/apply' or set using TF_VAR_MONGODB_USER=#{value}"
}

variable "region" {
  description = "Region to deploy AWS resources to"
}

variable "secret_key" {
  description = "AWS secret key"
}
