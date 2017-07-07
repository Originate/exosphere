resource "aws_ecs_task_definition" "task" {
  family = "${var.name}"

  container_definitions = <<EOF
[{
  "name": "${var.name}",
  "image": "${var.docker_image}",
  "command": ${jsonencode(var.command)},
  "cpu": ${var.cpu},
  "memory": ${var.memory},
  "portMappings": [
    ${var.container_port == "" ? "" : format("{\"containerPort\": %s}", var.container_port)}
  ],
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
  name              = "services/${var.env}/${var.name}"
  retention_in_days = 30
}
