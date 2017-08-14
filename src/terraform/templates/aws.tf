terraform {
  required_version = "= 0.9.11"

  backend "s3" {
    bucket     = "{{stateBucket}}"
    key        = "dev/terraform.tfstate"
    region     = "{{region}}"
    lock_table = "{{lockTable}}"
  }
}

provider "aws" {
  region              = "${var.region}"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["${var.account_id}"]
}

module "aws" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws?ref=8786f912"

  name     = "{{appName}}"
  env      = "production"
  key_name = "${var.key_name}"
}
