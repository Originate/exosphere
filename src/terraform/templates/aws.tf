variable "aws_profile" {
  default = "default"
}

terraform {
  required_version = ">= 0.10.0"

  backend "s3" {
    bucket         = "{{stateBucket}}"
    key            = "terraform.tfstate"
    region         = "{{region}}"
    dynamodb_table = "{{lockTable}}"
  }
}

provider "aws" {
  version = "0.1.4"

  region              = "{{region}}"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["{{accountID}}"]
}

variable "key_name" {
  default = ""
}

module "aws" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws?ref={{terraformCommitHash}}"

  name              = "{{appName}}"
  env               = "production"
  external_dns_name = "{{{url}}}"
  key_name          = "${var.key_name}"
}
