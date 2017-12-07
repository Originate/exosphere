variable "{{serviceRole}}_env_vars" {
  default = "[]"
}

variable "{{serviceRole}}_docker_image" {}

module "{{serviceRole}}" {
  source = "github.com/Originate/exosphere.git//terraform//aws//worker-service?ref={{terraformCommitHash}}"

  name = "{{serviceRole}}"

  cluster_id            = "${module.aws.ecs_cluster_id}"
  cpu                   = "{{cpu}}"
  desired_count         = 1
  docker_image          = "${var.{{serviceRole}}_docker_image}"
  env                   = "production"
  environment_variables = "${var.{{serviceRole}}_env_vars}"
  memory_reservation    = "{{memory}}"
  region                = "${module.aws.region}"
}
