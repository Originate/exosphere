module "aws" {
  source = "./aws"

  account_id      = "518695917306"
  env             = "production"
  key_name        = "hugo"
  region          = "us-west-2"
  security_groups = ["${module.exocom_public_service.security_groups}"]
}

module "exocom_public_service" {
  source = "./aws/service"

  name = "exocom"

  cluster_id            = "${module.aws.cluster_id}"
  command               = "bin/exocom"
  container_port        = "3100"
  cpu_units             = "128"
  docker_image          = "518695917306.dkr.ecr.us-west-2.amazonaws.com/exocom:0.22.1"
  ecs_role_arn          = "${module.aws.ecs_iam_role_arn}"
  env                   = "production"
  environment_variables = {}
  health_check_endpoint = "/config.json"
  memory_reservation    = "128"
  region                = "us-west-2"
  subnet_ids            = ["${module.aws.public_subnet_ids}"]
  vpc_id                = "${module.aws.vpc_id}"
}
