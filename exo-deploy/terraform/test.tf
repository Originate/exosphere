variable "MONGODB_USER" {
  description = "Environment variable for mlabs username. Prompted for during 'terraform plan/apply' or set using TF_VAR_MONGODB_USER=#{value}"
}

variable "MONGODB_PW" {
  description = "Environment variable for mlabs password. Prompted for during 'terraform plan/apply' or set using TF_VAR_MONGODB_PW=#{value}"
}

module "aws" {
  source = "./aws"

  account_id      = "518695917306"
  env             = "production"
  key_name        = "hugo"
  region          = "us-west-2"
  security_groups = ["${module.exocom.security_groups}"]
}

module "exocom" {
  source = "./aws/service"

  name = "exocom"

  alb_subnet_ids        = ["${module.aws.public_subnet_ids}"] //TODO: conditional?
  cluster_id            = "${module.aws.cluster_id}"
  command               = "bin/exocom" //TODO: make list
  container_port        = "3100"
  cpu_units             = "128"
  docker_image          = "518695917306.dkr.ecr.us-west-2.amazonaws.com/exocom:0.22.1"
  ecs_role_arn          = "${module.aws.ecs_iam_role_arn}"
  env                   = "production"
  environment_variables = {}
  health_check_endpoint = "/config.json"
  memory_reservation    = "128"
  region                = "us-west-2"
  vpc_id                = "${module.aws.vpc_id}"
}

module "web" {
  source = "./aws/service"

  name = "web"

  alb_subnet_ids        = ["${module.aws.public_subnet_ids}"]
  cluster_id            = "${module.aws.cluster_id}"
  command               = "node_modules/.bin/lsc app"
  container_port        = "3000"
  cpu_units             = "128"
  docker_image          = "518695917306.dkr.ecr.us-west-2.amazonaws.com/space-tweet-web-service:latest"
  ecs_role_arn          = "${module.aws.ecs_iam_role_arn}"
  env                   = "production"
  environment_variables = {
    ROLE        = "web"
    EXOCOM_HOST = "${module.exocom.url}"
    EXOCOM_PORT = "80"
  }
  health_check_endpoint = "/"
  memory_reservation    = "128"
  region                = "us-west-2"
  vpc_id                = "${module.aws.vpc_id}"
}

module "users" {
  source = "./aws/service"

  name = "users"

  alb_subnet_ids        = ["${module.aws.private_subnet_ids}"]
  cluster_id            = "${module.aws.cluster_id}"
  command               = "node_modules/exoservice/bin/exo-js"
  /* container_port        = "3000" */
  cpu_units             = "128"
  docker_image          = "518695917306.dkr.ecr.us-west-2.amazonaws.com/space-tweet-users-service:latest"
  ecs_role_arn          = "${module.aws.ecs_iam_role_arn}"
  env                   = "production"
  environment_variables = {
    ROLE         = "users"
    EXOCOM_HOST  = "${module.exocom.url}"
    EXOCOM_PORT  = "80"
    MONGODB_USER = "${var.MONGODB_USER}"
    MONGODB_PW   = "${var.MONGODB_PW}"
  }
  health_check_endpoint = "/"
  memory_reservation    = "128"
  region                = "us-west-2"
  vpc_id                = "${module.aws.vpc_id}"
}

module "tweets" {
  source = "./aws/service"

  name = "tweets"

  alb_subnet_ids        = ["${module.aws.private_subnet_ids}"]
  cluster_id            = "${module.aws.cluster_id}"
  command               = "node_modules/exoservice/bin/exo-js"
  /* container_port        = "3000" */
  cpu_units             = "128"
  docker_image          = "518695917306.dkr.ecr.us-west-2.amazonaws.com/space-tweet-tweets-service:latest"
  ecs_role_arn          = "${module.aws.ecs_iam_role_arn}"
  env                   = "production"
  environment_variables = {
    ROLE         = "tweets"
    EXOCOM_HOST  = "${module.exocom.url}"
    EXOCOM_PORT  = "80"
    MONGODB_USER = "${var.MONGODB_USER}"
    MONGODB_PW   = "${var.MONGODB_PW}"
  }
  health_check_endpoint = "/"
  memory_reservation    = "128"
  region                = "us-west-2"
  vpc_id                = "${module.aws.vpc_id}"
}
