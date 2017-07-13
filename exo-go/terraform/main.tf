terraform {
  required_version = "= 0.9.9"

  backend "s3" {
    bucket     = "space-tweet-terraform"
    key        = "dev/terraform.tfstate"
    region     = "us-west-2"
    lock_table = "TerraformLocks"
  }
}

module "aws" {
  source = "./aws"

  account_id       = "${var.account_id}"
  application_name = "space-tweet"
  env              = "production"
  key_name         = "${var.key_name}"
  region           = "${var.region}"
  security_groups  = []
  "asdlf;kj"
}
