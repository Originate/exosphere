variable "aws_profile" {
  default = "default"
}

variable "aws_region" {}

variable "aws_account_id" {}

variable "aws_ssl_certificate_arn" {}

variable "application_url" {}

variable "env" {}

terraform {
  required_version = "= {{{terraformVersion}}}"

  backend "s3" {
    key            = "terraform.tfstate"
    dynamodb_table = "{{lockTable}}"
  }
}

provider "aws" {
  version = "0.1.4"

  region              = "${var.aws_region}"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["${var.aws_account_id}"]
}

variable "key_name" {
  default = ""
}

module "aws" {
  source = "github.com/Originate/exosphere.git//terraform//aws?ref={{terraformCommitHash}}"

  name              = "{{appName}}"
  env               = "${var.env}"
  external_dns_name = "${var.application_url}"
  key_name          = "${var.key_name}"
  log_bucket_prefix = "${var.aws_account_id}-{{appName}}-${var.env}"
}
