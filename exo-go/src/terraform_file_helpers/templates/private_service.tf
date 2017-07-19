module "{{serviceRole}}" {
  source = "./aws/worker-service"

  name = "{{serviceRole}}"

  cluster_id            = "${module.aws.cluster_id}"
  command               = {{{startupCommand}}}
  cpu_units             = "{{cpu}}"
  env                   = "production"
  memory_reservation    = "{{memory}}"
  region                = "${var.region}"
}
