module "{{serviceRole}}" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//worker-service?ref=8786f912"

  name = "{{serviceRole}}"

  cluster_id    = "${module.aws.ecs_cluster_id}"
  command       = {{{startupCommand}}}
  cpu           = "{{cpu}}"
  desired_count = 1
  docker_image  = "{{{dockerImage}}}"
  env           = "production"
  memory        = "{{memory}}"
  region        = "${var.region}"
}
