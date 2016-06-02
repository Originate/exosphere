variable "name" {}
variable "cluster_id" {}
variable "security_groups" {}
variable "subnet_ids" {}


resource "aws_ecs_task_definition" "task" {
  family = "${var.name}-task-definition"
  container_definitions = "${file("${path.module}/wordpress.json")}"
}


resource "aws_iam_role" "ecs_service_role" {
  name = "${var.name}-role"
  assume_role_policy = "${file("${path.module}/ecs-role.json")}"
}


resource "aws_iam_role_policy" "ecs_service_role_policy" {
  name = "${var.name}-role-policy"
  role = "${aws_iam_role.ecs_service_role.id}"
  policy = "${file("${path.module}/ecs-role-policy.json")}"
}


resource "aws_ecs_service" "service" {
  name = "${var.name}"
  cluster = "${var.cluster_id}"
  task_definition = "${aws_ecs_task_definition.task.arn}"
  desired_count = 3
  deployment_minimum_healthy_percent = 100
  iam_role = "${aws_iam_role.ecs_service_role.arn}"
  depends_on = ["aws_iam_role_policy.ecs_service_role_policy"]

  load_balancer {
    elb_name = "${aws_elb.elb.id}"
    container_name = "wordpress"
    container_port = 80
  }
}


resource "aws_elb" "elb" {
  name = "${var.name}-elb"
  subnets = ["${split(",", var.subnet_ids)}"]
  security_groups = ["${var.security_groups}"]

  listener {
    instance_port = 80
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:80/"
    interval            = 30
  }
}
