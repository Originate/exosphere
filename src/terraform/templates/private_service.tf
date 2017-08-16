module "{{serviceRole}}" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//worker-service?ref={{terraformCommitHash}}"

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
