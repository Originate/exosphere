module "internal_alb" {
  source                = "../alb"

  env                   = "${var.env}"
  health_check_endpoint = "${var.health_check_endpoint}"
  internal              = true
  log_bucket            = "${var.log_bucket}"
  name                  = "${var.name}"
  security_groups       = ["${var.alb_security_group}"]
  subnet_ids            = "${var.alb_subnet_ids}"
  vpc_id                = "${var.vpc_id}"
}

module "task_definition" {
  source                = "../ecs-task-definition"

  command               = "${var.command}"
  container_port        = "${var.container_port}"
  cpu                   = "${var.cpu}"
  docker_image          = "${var.docker_image}"
  env                   = "${var.env}"
  environment_variables = "${var.environment_variables}"
  memory                = "${var.memory}"
  name                  = "${var.name}"
  region                = "${var.region}"
}

resource "aws_ecs_service" "service" {
  name                               = "${var.name}"
  cluster                            = "${var.cluster_id}"
  deployment_minimum_healthy_percent = 100
  desired_count                      = 1
  task_definition                    = "${module.task_definition.arn}"
  iam_role                           = "${var.ecs_role_arn}"

  depends_on                         = ["module.internal_alb"]

  load_balancer {
    container_name   = "${var.name}"
    container_port   = "${var.container_port}"
    target_group_arn = "${module.internal_alb.target_group_id}"
  }
}
