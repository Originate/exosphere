terraform {
  required_version = "= 0.9.9"

  backend "s3" {
    bucket     = "space-tweet-terraform"
    key        = "dev/terraform.tfstate"
    region     = "us-west-2"
    lock_table = "TerraformLocks"
  }
}

provider "aws" {
  region              = "${var.region}"
  profile             = "${var.aws_profile}"
  allowed_account_ids = ["${var.account_id}"]
}

module "aws" {
  source = "./aws"

  name     = "space-tweet"
  env      = "production"
  key_name = "${var.key_name}"
}

module "exocom_cluster" {
  source = "./aws/custom/exocom/exocom-cluster"

  availability_zones      = "${module.aws.availability_zones}"
  env                     = "production"
  internal_hosted_zone_id = "${module.aws.internal_zone_id}"
  instance_type           = "t2.micro"
  key_name                = "${var.key_name}"
  name                    = "exocom"
  region                  = "${module.aws.region}"

  bastion_security_group = ["${module.aws.bastion_security_group}"]

  ecs_cluster_security_groups = [ "${module.aws.ecs_cluster_security_group}",
    "${module.aws.external_alb_security_group}",
  ]

  subnet_ids = "${module.aws.private_subnet_ids}"
  vpc_id     = "${module.aws.vpc_id}"
}

module "exocom_service" {
  source = "./aws/custom/exocom/exocom-service"

  cluster_id     = "${module.exocom_cluster.cluster_id}"
  command        = ["bin/exocom"]
  container_port = "3100"
  cpu_units      = "128"
  docker_image   = "518695917306.dkr.ecr.us-west-2.amazonaws.com/exocom:latest"
  env            = "production"

  environment_variables = {
    DEBUG = "exocom,exocom:websocket-subsystem"
    ROLE  = "exocom"

    SERVICE_ROUTES = <<EOF
[{"role":"space-tweet-web-service","receives":["user details","user not found","user updated","user deleted","users listed","user created","tweets listed","tweet created","tweet deleted"],"sends":["get user details","delete user","update user","list users","create user","list tweets","create tweet","delete tweet"]},{"role":"exosphere-users-service","receives":["create user","list users","get user details","update user","delete user"],"sends":["user created","users listed","user details","user not found","user updated","user deleted","user not created"]},{"role":"exosphere-tweets-service","receives":["create tweet","list tweets","get tweet details","update tweet","delete tweet"],"sends":["tweet created","tweets listed","tweet details","tweet not found","tweet updated","tweet deleted","tweet not created"]}]
EOF
  }

  memory_reservation = "128"
  name               = "exocom"
  region             = "${module.aws.region}"
}

module "web" {
  source = "./aws/public-service"

  name               = "web"
  alb_security_group = "${module.aws.external_alb_security_group}"
  alb_subnet_ids     = ["${module.aws.public_subnet_ids}"]
  cluster_id         = "${module.aws.ecs_cluster_id}"
  command            = ["node_modules/.bin/lsc", "app"]
  container_port     = "3000"
  cpu                = "128"
  docker_image       = "518695917306.dkr.ecr.us-west-2.amazonaws.com/space-tweet-web-service:latest"
  ecs_role_arn       = "${module.aws.ecs_service_iam_role_arn}"
  env                = "production"

  environment_variables = {
    ROLE        = "space-tweet-web-service"
    EXOCOM_HOST = "${module.exocom_cluster.exocom_address}"
    EXOCOM_PORT = "80"
    DEBUG       = "exorelay,exorelay:message-manager"
  }

  external_dns_name     = "spacetweet.originate.com"
  external_zone_id      = "${var.hosted_zone_id}"
  health_check_endpoint = "/"
  internal_dns_name     = "spacetweet"
  internal_zone_id      = "${module.aws.internal_zone_id}"
  log_bucket            = "${module.aws.log_bucket_id}"
  memory                = "128"
  region                = "${module.aws.region}"
  ssl_certificate_arn   = "${var.ssl_certificate_arn}"
  vpc_id                = "${module.aws.vpc_id}"
}

module "users" {
  source = "./aws/worker-service"

  name         = "users"
  cluster_id   = "${module.aws.ecs_cluster_id}"
  command      = ["node_modules/exoservice/bin/exo-js"]
  cpu          = "128"
  docker_image = "518695917306.dkr.ecr.us-west-2.amazonaws.com/space-tweet-users-service:latest"
  env          = "production"

  environment_variables = {
    ROLE         = "exosphere-users-service"
    EXOCOM_HOST  = "${module.exocom_cluster.exocom_address}"
    EXOCOM_PORT  = "80"
    MONGODB_USER = "${var.mongodb_user}"
    MONGODB_PW   = "${var.mongodb_pw}"
    DEBUG        = "exorelay,exorelay:message-manager,exorelay:websocket-listener"
  }

  memory = "128"
  region = "${module.aws.region}"
}

module "tweets" {
  source = "./aws/worker-service"

  name         = "tweets"
  cluster_id   = "${module.aws.ecs_cluster_id}"
  command      = ["node_modules/exoservice/bin/exo-js"]
  cpu          = "128"
  docker_image = "518695917306.dkr.ecr.us-west-2.amazonaws.com/space-tweet-tweets-service:latest"
  env          = "production"

  environment_variables = {
    ROLE         = "exosphere-tweets-service"
    EXOCOM_HOST  = "${module.exocom_cluster.exocom_address}"
    EXOCOM_PORT  = "80"
    MONGODB_USER = "${var.mongodb_user}"
    MONGODB_PW   = "${var.mongodb_pw}"
    DEBUG        = "exorelay,exorelay:message-manager,exorelay:websocket-listener"
  }

  memory = "128"
  region = "${module.aws.region}"
}
