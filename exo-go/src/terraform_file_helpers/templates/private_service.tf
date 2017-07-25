module "{{serviceRole}}" {
  source = "./aws/worker-service"

  name = "{{serviceRole}}"

  cluster_id    = "${module.aws.cluster_id}"
  command       = {{{startupCommand}}}
  cpu           = "{{cpu}}"
  desired_count = 1
  env           = "production"
  memory        = "{{memory}}"
  region        = "${var.region}"
}
