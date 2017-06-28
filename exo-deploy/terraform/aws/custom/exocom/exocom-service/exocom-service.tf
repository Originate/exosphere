module "task_definition" {
  source                = "./exocom-task-definition"

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
  task_definition                    = "${module.task_definition.task_arn}"
}
