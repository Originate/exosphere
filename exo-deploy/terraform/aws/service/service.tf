module "iam" {
  source = "./iam"

  env = "${var.env}"
}

resource "aws_ecs_service" "service" {
  name                               = "${var.name}"
  cluster                            = "${var.cluster_id}"
  task_definition                    = "${aws_ecs_task_definition.task.arn}"
  desired_count                      = 1
  deployment_minimum_healthy_percent = 100
  iam_role                           = "${module.iam.iam_role_arn}"
  depends_on                         = ["aws_alb.alb"]

  load_balancer {
    target_group_arn = "${aws_alb_target_group.target_group.id}"
    container_name   = "${var.name}"
    container_port   = "${var.container_port}"
  }
}

resource "aws_ecs_task_definition" "task" {
  family = "${var.name}"

  container_definitions = <<EOF
[{
  "name": "${var.name}",
  "image": "${var.docker_image}",
  "command": ["${var.command}"],
  "cpu": ${var.cpu_units},
  "memoryReservation": ${var.memory_reservation},
  "portMappings": [{
    "containerPort": ${var.container_port}
  }],
  "environment": [
    ${join(",",
           formatlist("{\"name\": %q, \"value\": %q}",
                      keys(var.environment_variables),
                      values(var.environment_variables)
     ))}
  ],
  "logConfiguration": {
    "logDriver": "awslogs",
    "options": {
      "awslogs-region": "${var.region}",
      "awslogs-group": "${aws_cloudwatch_log_group.log_group.name}"
    }
  },
  "essential": true
}]
EOF
}

resource "aws_cloudwatch_log_group" "log_group" {
  name = "services/${var.env}/${var.name}"
}
