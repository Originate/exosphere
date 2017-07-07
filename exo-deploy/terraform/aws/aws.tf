data "aws_availability_zones" "available" {}

module "internal_dns" {
  source = "./internal-dns"

  name    = "${var.name}.local"
  vpc_id  = "${module.network.vpc_id}"
  servers = ["${cidrhost(module.network.vpc_cidr, 2)}"]
}

module "network" {
  source = "./network"

  name               = "${var.name}-${var.env}"
  env                = "${var.env}"
  availability_zones = "${data.aws_availability_zones.available.names}"
  region             = "${var.region}"
  key_name           = "${var.key_name}"
}

module "alb_security_groups" {
  source = "./alb-security-groups"

  name     = "${var.name}-${var.env}"
  env      = "${var.env}"
  vpc_cidr = "${module.network.vpc_cidr}"
  vpc_id   = "${module.network.vpc_id}"
}

module "ecs_cluster" {
  source = "./ecs-cluster"

  name          = "${var.name}-${var.env}"
  env           = "${var.env}"
  region        = "${var.region}"
  instance_type = "${var.ecs_instance_type}"
  ebs_optimized = "${var.ecs_ebs_optimized}"
  key_name      = "${var.key_name}"

  alb_security_groups = ["${module.alb_security_groups.internal_alb_id}",
    "${module.alb_security_groups.external_alb_id}",
  ]

  bastion_security_group = "${module.network.bastion_security_group_id}"
  subnet_ids             = ["${module.network.private_subnet_ids}"]

  vpc_id = "${module.network.vpc_id}"
}

module "s3_logs" {
  source = "./s3-logs"

  name = "${var.name}-${var.env}"
  env  = "${var.env}"
}
