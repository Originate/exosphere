resource "aws_ecs_task_definition" "task" {
  family = "${var.name}"

  container_definitions = <<EOF
[{
  "name": "${var.name}",
  "image": "${var.docker_image}",
  "command": ${jsonencode(var.command)},
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

output "task_arn" {
  value = "${aws_ecs_task_definition.task.arn}"
  description = "ARN of task definition to be passed to ECS service"
}
