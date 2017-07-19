module "{{serviceRole}}" {
  source = "./aws/public-service"

  name = "{{serviceRole}}"

  alb_security_group      = ["${module.aws.external_alb_security_group}"]
  alb_subnet_ids          = ["${module.aws.public_subnet_ids}"]
  cluster_id              = "${module.aws.cluster_id}"
  command                 = {{{startupCommand}}}
  container_port          = "{{publicPort}}"
  cpu_units               = "{{cpu}}"
  ecs_role_arn            = "${module.aws.ecs_service_iam_role_arn}"
  env                     = "production"
  external_dns_name       = "{{{url}}}"
  external_hosted_zone_id = "${var.hosted_zone_id}"
  health_check_endpoint   = "{{{healthCheck}}}"
  internal_dns_name       = "${module.aws.internal_dns_name}"
  internal_hosted_zone_id = "${module.aws.internal_hosted_zone_id}"
  log_bucket              = "${module.aws.log_bucket_id}"
  memory_reservation      = "{{memory}}"
  region                  = "${var.region}"
  vpc_id                  = "${module.aws.vpc_id}"
}
