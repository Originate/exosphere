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
    key            = "services.tfstate"
    dynamodb_table = "{{lockTable}}"
  }
}

provider "aws" {
  version = "0.1.4"

  region              = "${var.aws_region}"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["${var.aws_account_id}"]
}

data "terraform_remote_state" "main_infrastructure" {
  backend = "s3"
  config {
    key            = "${var.aws_account_id}-{{appName}}-${var.env}-infrastructure-terraform/terraform.tfstate"
    dynamodb_table = "{{lockTable}}"
    region         = "${var.aws_region}"
    profile        = "${var.aws_profile}"
  }
}
