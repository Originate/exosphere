variable "{{serviceRole}}_env_vars" {
  default = "[]"
}

variable "{{serviceRole}}_docker_image" {}

module "{{serviceRole}}" {
  source = "github.com/Originate/exosphere.git//terraform//aws//private-service?ref={{terraformCommitHash}}"

  name = "{{serviceRole}}"

  alb_security_group    = "${data.terraform_remote_state.main_infrastructure.internal_alb_security_group}"
  alb_subnet_ids        = ["${data.terraform_remote_state.main_infrastructure.private_subnet_ids}"]
  cluster_id            = "${data.terraform_remote_state.main_infrastructure.ecs_cluster_id}"
  container_port        = "{{{publicPort}}}"
  cpu                   = "{{cpu}}"
  desired_count         = 1
  docker_image          = "${var.{{serviceRole}}_docker_image}"
  ecs_role_arn          = "${data.terraform_remote_state.main_infrastructure.ecs_service_iam_role_arn}"
  env                   = "${var.env}"
  environment_variables = "${var.{{serviceRole}}_env_vars}"
  health_check_endpoint = "{{{healthCheck}}}"
  internal_dns_name     = "{{{serviceRole}}}"
  internal_zone_id      = "${data.terraform_remote_state.main_infrastructure.internal_zone_id}"
  log_bucket            = "${data.terraform_remote_state.main_infrastructure.log_bucket_id}"
  memory_reservation    = "{{memory}}"
  region                = "${data.terraform_remote_state.main_infrastructure.region}"
  vpc_id                = "${data.terraform_remote_state.main_infrastructure.vpc_id}"
}
