data "aws_availability_zones" "available" {}

provider "aws" {
  region              = "${var.region}"
  allowed_account_ids = ["${var.account_id}"]
}

/* module "dhcp" { */
/*   source  = "./dhcp" */
/*   name    = "${module.dns.name}" */
/*   vpc_id  = "${module.network.vpc_id}" */
/*   servers = "${coalesce()}" //TODO */
/* } */
/*  */
/* module "dns" { */
/*   source = "./dns" */
/*   name   = "${var.domain_name}" */
/*   vpc_id = "${module.network.vpc_id}" */
/* } */

module "network" {
  source             = "./network"

  env                = "${var.env}"
  availability_zones = "${data.aws_availability_zones.available.names}"
  region             = "${var.region}"
  key_name           = "${var.key_name}"
}

module "alb_security_groups" {
  source = "./alb-security-groups"

  name   = "${var.application_name}"
  env    = "${var.env}"
  vpc_id = "${module.network.vpc_id}"
}

module "cluster" {
  source               = "./cluster"

  availability_zones = "${data.aws_availability_zones.available.names}"
  env                = "${var.env}"
  instance_type      = "t2.micro"
  key_name           = "${var.key_name}"
  name               = "exosphere-cluster"
  region             = "${var.region}"
  security_groups    = ["${module.network.bastion_security_group_id}",
                        "${module.alb_security_groups.internal_alb_security_group}",
                        "${module.alb_security_groups.external_alb_security_group}",
                        "${var.security_groups}"]
  subnet_ids         = ["${module.network.private_subnet_ids}"]
  vpc_id             = "${module.network.vpc_id}"
}

/* module "s3_logs" { */
/*   source                  = "./s3-logs" */
/*   name                    = "${var.name}" */
/*   env                     = "${var.env}" */
/*   account_id              = "${var.account_id}" */
/* } */
