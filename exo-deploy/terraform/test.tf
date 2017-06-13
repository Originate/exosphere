module "aws" {
  source = "./aws"

  env = "production"
  region = "us-west-2"
  account_id = "518695917306"
  security_groups = ["${module.exocom.security_groups}"]
}

module "exocom" {
  source = "./aws/service"

  region = "us-west-2"
  env = "production"
  name = "exocom"
  command = "bin/exocom"
  vpc_id = "${module.aws.vpc_id}"
  public_subnet_ids = ["${module.aws.public_subnet_ids}"]
  cluster_id = "${module.aws.cluster_id}"
  cpu_units = "128"
  memory_reservation = "128"
  docker_image = "518695917306.dkr.ecr.us-west-2.amazonaws.com/exocom:0.22.1"
  container_port = "3100"
  environment_variables = {}
}
