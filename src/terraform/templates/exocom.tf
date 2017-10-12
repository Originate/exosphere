module "exocom_cluster" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//dependencies//exocom//exocom-cluster?ref={{terraformCommitHash}}"

  availability_zones          = "${module.aws.availability_zones}"
  env                         = "production"
  internal_hosted_zone_id     = "${module.aws.internal_zone_id}"
  instance_type               = "t2.micro"
  key_name                    = "${var.key_name}"
  name                        = "exocom"
  region                      = "${module.aws.region}"

  bastion_security_group      = ["${module.aws.bastion_security_group}"]

  ecs_cluster_security_groups = [ "${module.aws.ecs_cluster_security_group}",
    "${module.aws.external_alb_security_group}",
  ]

  subnet_ids                  = "${module.aws.private_subnet_ids}"
  vpc_id                      = "${module.aws.vpc_id}"
}

variable "exocom_service_routes" {
  default = ""
}

module "exocom_service" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//dependencies//exocom//exocom-service?ref={{terraformCommitHash}}"

  cluster_id            = "${module.exocom_cluster.cluster_id}"
  cpu_units             = "128"
  docker_image          = "{{{dockerImage}}}"
  env                   = "production"
  environment_variables = {
    ROLE           = "exocom"
    SERVICE_ROUTES = "${var.exocom_service_routes}"
  }
  memory_reservation    = "128"
  name                  = "exocom"
  region                = "${module.aws.region}"
}
