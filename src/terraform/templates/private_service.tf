module "{{serviceRole}}" {
  source = "./aws/worker-service"

  name = "{{serviceRole}}"

  cluster_id    = "${module.aws.cluster_id}"
  command       = {{{startupCommand}}}
  cpu           = "{{cpu}}"
  desired_count = 1
  docker_image  = "{{{dockerImage}}}"
  env           = "production"
  memory        = "{{memory}}"
  region        = "${var.region}"
}
