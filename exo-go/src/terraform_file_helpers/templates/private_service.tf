module "{{serviceRole}}" {
  source = "./aws/worker-service"

  name = "{{serviceRole}}"

  cluster_id            = "${module.aws.cluster_id}"
  command               = {{{startupCommand}}}
  cpu_units             = "{{cpu}}"
  docker_image          = "" //TODO: implement after ecr functionality is in place
  env                   = "production"
  environment_variables = {
    {{{envVars}}}
  }
  memory_reservation    = "{{memory}}"
  region                = "${var.region}"
}
