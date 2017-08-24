variable "region" {}
variable "aws_profile" {}
variable "account_id" {}

terraform {
  required_version = ">= 0.10.0"

  backend "s3" {
    bucket         = "{{stateBucket}}"
    key            = "dev/terraform.tfstate"
    region         = "{{region}}"
    dynamodb_table = "{{lockTable}}"
  }
}

provider "aws" {
  version = "0.1.4"

  region              = "${var.region}"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["${var.account_id}"]
}

module "aws" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws?ref={{terraformCommitHash}}"

  name              = "{{appName}}"
  env               = "production"
  external_dns_name = "{{{url}}}"
  key_name          = "${var.key_name}"
}
