variable "name" {}
variable "cluster_id" {}
variable "security_groups" {}
variable "subnet_ids" {}
variable "service_type" {}


resource "aws_ecs_task_definition" "task" {
  family = "${var.name}-task-definition"
  container_definitions = "${file("${path.root}/${var.service_type}-container-definition.json")}"
}


resource "aws_ecs_service" "service" {
  name = "${var.name}"
  cluster = "${var.cluster_id}"
  task_definition = "${aws_ecs_task_definition.task.arn}"
  desired_count = 1
  deployment_minimum_healthy_percent = 100
}
