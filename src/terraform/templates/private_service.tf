variable "{{serviceRole}}_env_vars" {
  default = "[]"
}

module "{{serviceRole}}" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//private-service?ref={{terraformCommitHash}}"

  name = "{{serviceRole}}"

  alb_security_group    = "${module.aws.internal_alb_security_group}"
  alb_subnet_ids        = ["${module.aws.private_subnet_ids}"]
  cluster_id            = "${module.aws.ecs_cluster_id}"
  command               = {{{startupCommand}}}
  container_port        = "{{{publicPort}}}"
  cpu                   = "{{cpu}}"
  desired_count         = 1
  docker_image          = "{{{dockerImage}}}"
  ecs_role_arn          = "${module.aws.ecs_service_iam_role_arn}"
  env                   = "production"
  environment_variables = "${var.{{serviceRole}}_env_vars}"
  health_check_endpoint = "{{{healthCheck}}}"
  internal_dns_name     = "{{{serviceRole}}}"
  internal_zone_id      = "${module.aws.internal_zone_id}"
  log_bucket            = "${module.aws.log_bucket_id}"
  memory                = "{{memory}}"
  region                = "${var.region}"
  vpc_id                = "${module.aws.vpc_id}"
}
