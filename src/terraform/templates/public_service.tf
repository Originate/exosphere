variable "{{serviceRole}}_env_vars" {
  default = "[]"
}

variable "{{serviceRole}}_docker_image" {}

module "{{serviceRole}}" {
  source = "git@github.com:Originate/exosphere.git//terraform//aws//public-service?ref={{terraformCommitHash}}"

  name = "{{serviceRole}}"

  alb_security_group    = "${module.aws.external_alb_security_group}"
  alb_subnet_ids        = ["${module.aws.public_subnet_ids}"]
  cluster_id            = "${module.aws.ecs_cluster_id}"
  container_port        = "{{publicPort}}"
  cpu                   = "{{cpu}}"
  desired_count         = 1
  docker_image          = "${var.{{serviceRole}}_docker_image}"
  ecs_role_arn          = "${module.aws.ecs_service_iam_role_arn}"
  env                   = "production"
  environment_variables = "${var.{{serviceRole}}_env_vars}"
  external_dns_name     = "{{{url}}}"
  external_zone_id      = "${module.aws.external_zone_id}"
  health_check_endpoint = "{{{healthCheck}}}"
  internal_dns_name     = "{{{serviceRole}}}"
  internal_zone_id      = "${module.aws.internal_zone_id}"
  log_bucket            = "${module.aws.log_bucket_id}"
  memory_reservation    = "{{memory}}"
  region                = "${module.aws.region}"
  ssl_certificate_arn   = "{{{sslCertificateArn}}}"
  vpc_id                = "${module.aws.vpc_id}"
}
