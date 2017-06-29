module "internal_alb" {
  source                = "../alb"

  env                   = "${var.env}"
  health_check_endpoint = "${var.health_check_endpoint}"
  log_bucket            = "${var.log_bucket}"
  name                  = "${var.name}"
  security_group        = "${var.alb_security_group}"
  subnet_ids            = "${var.alb_subnet_ids}"
  vpc_id                = "${var.vpc_id}"
}

module "task_definition" {
  source                = "../ecs-task-definition"

  command               = "${var.command}"
  container_port        = "${var.container_port}"
  cpu_units             = "${var.cpu_units}"
  docker_image          = "${var.docker_image}"
  env                   = "${var.env}"
  environment_variables = "${var.environment_variables}"
  memory_reservation    = "${var.memory_reservation}"
  name                  = "${var.name}"
  region                = "${var.region}"
}

resource "aws_ecs_service" "service" {
  name                               = "${var.name}"

  cluster                            = "${var.cluster_id}"
  deployment_minimum_healthy_percent = 100
  desired_count                      = 1
  iam_role                           = "${var.ecs_role_arn}"
  task_definition                    = "${module.task_definition.task_arn}"

  load_balancer {
    container_name   = "${var.name}"
    container_port   = "${var.container_port}"
    target_group_arn = "${module.internal_alb.target_group_id}"
  }
}
