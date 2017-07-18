module "exocom_cluster" {
  source = "./aws/custom/exocom/exocom-cluster"

  availability_zones = "${module.aws.availability_zones}"
  env                = "production"
  domain_name        = "{{{dnsName}}}"
  hosted_zone_id     = "${var.hosted_zone_id}"
  instance_type      = "t2.micro"
  key_name           = "${var.key_name}"
  name               = "exocom"
  region             = "${module.aws.region}"

  security_groups = ["${module.aws.bastion_security_group}",
    "${module.aws.ecs_cluster_security_group}",
    "${module.aws.external_alb_security_group}",
  ]

  subnet_ids = "${module.aws.private_subnet_ids}"
  vpc_id     = "${module.aws.vpc_id}"
}

module "exocom_service" {
  source = "./aws/custom/exocom/exocom-service"

  cluster_id     = "${module.exocom_cluster.cluster_id}"
  command        = ["bin/exocom"]
  container_port = "3100"
  cpu_units      = "128"
  docker_image   = "{{version}}" //TODO: implement once ecr support is in place
  env            = "production"
  environment_variables = {
    ROLE  = "exocom"
    SERVICE_ROUTES = <<EOF
{{{serviceRoutes}}}
EOF
  }

  memory_reservation = "128"
  name               = "exocom"
  region             = "${module.aws.region}"
}
