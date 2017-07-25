module "internal_alb" {
  source = "../alb"

  env                   = "${var.env}"
  health_check_endpoint = "${var.health_check_endpoint}"
  internal              = true
  internal_dns_name     = "${var.internal_dns_name}"
  internal_zone_id      = "${var.internal_zone_id}"
  log_bucket            = "${var.log_bucket}"
  name                  = "${var.env}-${var.name}"
  security_groups       = ["${var.alb_security_group}"]
  subnet_ids            = "${var.alb_subnet_ids}"
  vpc_id                = "${var.vpc_id}"
}

module "task_definition" {
  source = "../ecs-task-definition"

  command               = "${var.command}"
  container_port        = "${var.container_port}"
  cpu                   = "${var.cpu}"
  docker_image          = "${var.docker_image}"
  env                   = "${var.env}"
  environment_variables = "${var.environment_variables}"
  memory                = "${var.memory}"
  name                  = "${var.env}-${var.name}"
  region                = "${var.region}"
}

resource "aws_ecs_service" "service" {
  name                               = "${var.name}"
  cluster                            = "${var.cluster_id}"
  deployment_minimum_healthy_percent = 100
  desired_count                      = "${var.desired_count}"
  task_definition                    = "${module.task_definition.arn}"
  iam_role                           = "${var.ecs_role_arn}"

  load_balancer {
    container_name   = "${var.env}-${var.name}"
    container_port   = "${var.container_port}"
    target_group_arn = "${module.internal_alb.target_group_id}"
  }
}
