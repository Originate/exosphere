variable "name" {}
variable "cluster_id" {}


/* ECS task definitions are used to start new tasks (Docker containers)
on ECS clusters. container_definitions points to a JSON file that defines the
options to start a container with (similar to the flags that are passed
to a Docker container on startup). */
resource "aws_ecs_task_definition" "task" {
  family = "${var.name}-task-definition"
  container_definitions = "${file("${path.root}/exocom-container-definition.json")}"
}


resource "aws_ecs_service" "service" {
  name = "${var.name}"
  cluster = "${var.cluster_id}"
  task_definition = "${aws_ecs_task_definition.task.arn}"
  desired_count = 1
  deployment_minimum_healthy_percent = 100
}
