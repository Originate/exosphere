terraform {
  required_version = "= 0.9.9"

  backend "s3" {
    bucket     = "{{appName}}-terraform"
    key        = "dev/terraform.tfstate"
    region     = "{{region}}"
    lock_table = "TerraformLocks"
  }
}

module "aws" {
  source = "./aws"

  account_id       = "${var.account_id}"
  application_name = "{{appName}}"
  env              = "production"
  key_name         = "${var.key_name}"
  region           = "${var.region}"
  security_groups  = []
}
