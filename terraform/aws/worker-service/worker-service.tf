module "task_definition" {
  source = "../ecs-task-definition"

  command               = "${var.command}"
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
  iam_role                           = "${aws_iam_role.allow_dns_changes.arn}"
}

resource "aws_iam_role" "allow_dns_changes" {
  name = "ecs-allow-dns-changes"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "route53:*"
      ],
			"Principal": {
			  "Service": ["ecs.amazonaws.com"]
			}
    }
  ]
}
EOF
}
