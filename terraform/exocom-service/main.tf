variable "name" {}
variable "cluster_id" {}


resource "aws_ecs_task_definition" "task" {
  family = "${var.name}-task-definition"
  container_definitions = "${file("${path.module}/exocom.json")}"
}


resource "aws_ecs_service" "service" {
  name = "${var.name}"
  cluster = "${var.cluster_id}"
  task_definition = "${aws_ecs_task_definition.task.arn}"
  desired_count = 1
  deployment_minimum_healthy_percent = 100
}
