module "exocom_cluster" {
  source = "github.com/Originate/exosphere.git//remote-dependency-templates//exocom//modules//exocom-cluster?ref={{terraformCommitHash}}"

  availability_zones      = "${module.aws.availability_zones}"
  env                     = "${var.env}"
  internal_hosted_zone_id = "${module.aws.internal_zone_id}"
  instance_type           = "t2.micro"
  key_name                = "${var.key_name}"
  name                    = "exocom"
  region                  = "${module.aws.region}"

  bastion_security_group = ["${module.aws.bastion_security_group}"]

  ecs_cluster_security_groups = ["${module.aws.ecs_cluster_security_group}",
    "${module.aws.external_alb_security_group}",
  ]

  subnet_ids = "${module.aws.private_subnet_ids}"
  vpc_id     = "${module.aws.vpc_id}"
}

variable "exocom_env_vars" {
  default = ""
}

module "exocom_service" {
  source = "github.com/Originate/exosphere.git//remote-dependency-templates//exocom//modules//exocom-service?ref={{terraformCommitHash}}"

  cluster_id            = "${module.exocom_cluster.cluster_id}"
  cpu_units             = "128"
  docker_image          = "originate/exocom:{{version}}"
  env                   = "${var.env}"
  environment_variables = "${var.exocom_env_vars}"
  memory_reservation    = "128"
  name                  = "exocom"
  region                = "${module.aws.region}"
}
