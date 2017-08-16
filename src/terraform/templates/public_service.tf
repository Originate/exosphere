module "{{serviceRole}}" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//public-service?ref={{terraformCommitHash}}"

  name = "{{serviceRole}}"

  alb_security_group    = ["${module.aws.external_alb_security_group}"]
  alb_subnet_ids        = ["${module.aws.public_subnet_ids}"]
  cluster_id            = "${module.aws.cluster_id}"
  command               = {{{startupCommand}}}
  container_port        = "{{publicPort}}"
  cpu                   = "{{cpu}}"
  desired_count         = 1
  docker_image          = "{{{dockerImage}}}"
  ecs_role_arn          = "${module.aws.ecs_service_iam_role_arn}"
  env                   = "production"
  external_dns_name     = "{{{url}}}"
  external_zone_id      = "${var.hosted_zone_id}"
  health_check_endpoint = "{{{healthCheck}}}"
  internal_dns_name     = "${module.aws.internal_dns_name}"
  internal_zone_id      = "${module.aws.internal_hosted_zone_id}"
  log_bucket            = "${module.aws.log_bucket_id}"
  memory                = "{{memory}}"
  region                = "${var.region}"
  ssl_certificate_arn   = "${var.ssl_certificate_arn}"
  vpc_id                = "${module.aws.vpc_id}"
}
