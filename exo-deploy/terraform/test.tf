variable "MONGODB_USER" {
  description = "Environment variable for mlabs username. Prompted for during 'terraform plan/apply' or set using TF_VAR_MONGODB_USER=#{value}"
}

variable "MONGODB_PW" {
  description = "Environment variable for mlabs password. Prompted for during 'terraform plan/apply' or set using TF_VAR_MONGODB_PW=#{value}"
}

module "aws" {
  source = "./aws"

  account_id       = "518695917306"
  application_name = "space tweet"
  env              = "production"
  key_name         = "hugo"
  region           = "us-west-2"
  security_groups  = []
}

module "exocom_cluster" {
  source             = "./aws/custom/exocom/exocom-cluster"

  availability_zones = "${module.aws.availability_zones}"
  env                = "production"
  instance_type      = "t2.micro"
  key_name           = "hugo"
  name               = "exocom"
  region             = "us-west-2"
  security_groups    = ["${module.aws.bastion_security_group_id}", "${module.aws.cluster_security_group}"]
  subnet_ids         = "${module.aws.private_subnet_ids}"
  vpc_id             = "${module.aws.vpc_id}"
}

module "exocom_service" {
  source                      = "./aws/custom/exocom/exocom-service"


  cluster_id                  = "${module.exocom_cluster.cluster_id}"
  command                     = ["bin/exocom"]
  container_port              = "3100"
  cpu_units                   = "128"
  docker_image                = "518695917306.dkr.ecr.us-west-2.amazonaws.com/exocom:latest"
  ecs_role_arn                = "${module.aws.ecs_service_iam_role_arn}"
  elb_subnet_ids              = ["${module.aws.public_subnet_ids}"]
  env                         = "production"
  environment_variables       = {
    DEBUG = "exocom,exocom:websocket-subsystem"
    ROLE = "exocom"
    SERVICE_ROUTES = <<EOF
[{"role":"space-tweet-web-service","receives":["users.listed","users.created"],"sends":["users.list","users.create"]},{"role":"exosphere-users-service","receives":["users.create","users.list","user.get-details","user.update","user.delete"],"sends":["users.created","users.listed","user.details","user.not-found","user.updated","user.deleted","users.not-created"]},{"role":"exosphere-tweets-service","receives":["tweets.create","tweets.list","tweets.get-details","tweets.update","tweets.delete"],"sends":["tweets.created","tweets.listed","tweets.details","tweets.not-found","tweets.updated","tweets.deleted","tweets.not-created"]}]
EOF
  }
  health_check_endpoint       = "/config.json"
  memory_reservation          = "128"
  name                        = "exocom"
  region                      = "us-west-2"
  security_groups             = ["${module.exocom_cluster.security_groups}"]
  vpc_id                      = "${module.aws.vpc_id}"
}

module "web" {
  source = "./aws/public-service"

  name = "web"

  alb_security_group    = ["${module.aws.external_alb_security_group}"]
  alb_subnet_ids        = ["${module.aws.public_subnet_ids}"]
  cluster_id            = "${module.aws.cluster_id}"
  command               = ["node_modules/.bin/lsc", "app"]
  container_port        = "3000"
  cpu_units             = "128"
  docker_image          = "518695917306.dkr.ecr.us-west-2.amazonaws.com/space-tweet-web-service:latest"
  ecs_role_arn          = "${module.aws.ecs_service_iam_role_arn}"
  env                   = "production"
  environment_variables = {
    ROLE        = "space-tweet-web-service"
    EXOCOM_HOST = "${module.exocom_service.url}"
    EXOCOM_PORT = "80"
  }
  health_check_endpoint = "/health-check"
  memory_reservation    = "128"
  region                = "us-west-2"
  vpc_id                = "${module.aws.vpc_id}"
}

module "users" {
  source = "./aws/worker-service"

  name = "users"

  cluster_id            = "${module.aws.cluster_id}"
  command               = ["node_modules/exoservice/bin/exo-js"]
  cpu_units             = "128"
  docker_image          = "518695917306.dkr.ecr.us-west-2.amazonaws.com/space-tweet-users-service:latest"
  env                   = "production"
  environment_variables = {
    ROLE         = "exosphere-users-service"
    EXOCOM_HOST  = "${module.exocom_service.url}"
    EXOCOM_PORT  = "80"
    MONGODB_USER = "${var.MONGODB_USER}"
    MONGODB_PW   = "${var.MONGODB_PW}"
  }
  memory_reservation    = "128"
  region                = "us-west-2"
}

module "tweets" {
  source = "./aws/worker-service"

  name = "tweets"

  cluster_id            = "${module.aws.cluster_id}"
  command               = ["node_modules/exoservice/bin/exo-js"]
  cpu_units             = "128"
  docker_image          = "518695917306.dkr.ecr.us-west-2.amazonaws.com/space-tweet-tweets-service:latest"
  env                   = "production"
  environment_variables = {
    ROLE         = "exosphere-tweets-service"
    EXOCOM_HOST  = "${module.exocom_service.url}"
    EXOCOM_PORT  = "80"
    MONGODB_USER = "${var.MONGODB_USER}"
    MONGODB_PW   = "${var.MONGODB_PW}"
  }
  memory_reservation    = "128"
  region                = "us-west-2"
}
