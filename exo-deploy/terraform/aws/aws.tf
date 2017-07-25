data "aws_region" "current" {
  current = true
}

data "aws_availability_zones" "available" {}

module "network" {
  source = "./network"

  name               = "${var.env}-${var.name}"
  env                = "${var.env}"
  availability_zones = "${slice(data.aws_availability_zones.available.names, 0, 3)}"
  region             = "${data.aws_region.current.name}"
  key_name           = "${var.key_name}"
}

module "internal_dns" {
  source = "./internal-dns"

  name    = "${var.name}.local"
  env     = "${var.env}"
  vpc_id  = "${module.network.vpc_id}"
  servers = ["${cidrhost(module.network.vpc_cidr, 2)}"]
}

module "alb_security_groups" {
  source = "./alb/security-groups"

  name     = "${var.env}-${var.name}"
  env      = "${var.env}"
  vpc_cidr = "${module.network.vpc_cidr}"
  vpc_id   = "${module.network.vpc_id}"
}

module "ecs_cluster" {
  source = "./ecs-cluster"

  name                   = "${var.env}-${var.name}"
  env                    = "${var.env}"
  region                 = "${data.aws_region.current.name}"
  instance_type          = "${var.ecs_instance_type}"
  ebs_optimized          = "${var.ecs_ebs_optimized}"
  key_name               = "${var.key_name}"

  alb_security_groups    = ["${module.alb_security_groups.internal_id}",
    "${module.alb_security_groups.external_id}",
  ]

  bastion_security_group = "${module.network.bastion_security_group}"
  subnet_ids             = ["${module.network.private_subnet_ids}"]

  vpc_id = "${module.network.vpc_id}"
}

module "s3_logs" {
  source = "./s3-logs"

  name = "${var.env}-${var.name}"
  env  = "${var.env}"
}
