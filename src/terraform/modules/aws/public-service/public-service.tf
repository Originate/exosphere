module "task_definition" {
  source = "../ecs-task-definition"

  command               = "${var.command}"
  container_port        = "${var.container_port}"
  cpu                   = "${var.cpu}"
  docker_image          = "${var.docker_image}"
  env                   = "${var.env}"
  environment_variables = "${var.environment_variables}"
  memory_reservation    = "${var.memory_reservation}"
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
    target_group_arn = "${aws_alb_target_group.target_group.id}"
  }

  depends_on = ["aws_alb.alb"]
}
