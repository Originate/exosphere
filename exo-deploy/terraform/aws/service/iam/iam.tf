variable "env" {}

resource "aws_iam_role" "ecs" {
  name = "${var.env}-ecs-service-role"

  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": [
          "ecs.amazonaws.com",
          "ec2.amazonaws.com"
        ]
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "ecs_service" {
  name = "${var.env}-ecs-service-role-policy"
  role = "${aws_iam_role.ecs.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:Describe*",
        "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
        "elasticloadbalancing:DeregisterTargets",
        "elasticloadbalancing:Describe*",
        "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
        "elasticloadbalancing:RegisterTargets"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

/* resource "aws_iam_role_policy" "ecs_instance" { */
/*   name = "${var.env}-ecs-instance-role-policy" */
/*   role = "${aws_iam_role.ecs.id}" */
/*  */
/*   policy = <<EOF */
/* { */
/*   "Version": "2012-10-17", */
/*   "Statement": [ */
/*     { */
/*       "Effect": "Allow", */
/*       "Action": [ */
/*         "ecs:CreateCluster", */
/*         "ecs:DeregisterContainerInstance", */
/*         "ecs:DiscoverPollEndpoint", */
/*         "ecs:Poll", */
/*         "ecs:RegisterContainerInstance", */
/*         "ecs:StartTelemetrySession", */
/*         "ecs:Submit*", */
/*         "ecr:GetAuthorizationToken", */
/*         "ecr:BatchCheckLayerAvailability", */
/*         "ecr:GetDownloadUrlForLayer", */
/*         "ecr:BatchGetImage", */
/*         "ecs:StartTask", */
/*         "autoscaling:*" */
/*       ], */
/*       "Resource": "*" */
/*     }, */
/*     { */
/*       "Effect": "Allow", */
/*       "Action": [ */
/*         "logs:CreateLogGroup", */
/*         "logs:CreateLogStream", */
/*         "logs:PutLogEvents", */
/*         "logs:DescribeLogStreams" */
/*       ], */
/*       "Resource": "arn:aws:logs:*:*:*" */
/*     } */
/*   ] */
/* } */
/* EOF */
/* } */

/* resource "aws_iam_instance_profile" "ecs" { */
/*   name = "${var.env}-ecs-instance-profile" */
/*   path = "/" */
/*   role = "${aws_iam_role.ecs.name}" */
/* } */
/*  */
/* output "iam_instance_profile" { */
/*   value = "${aws_iam_instance_profile.ecs.arn}" */
/* } */

output "iam_role_arn" {
  value = "${aws_iam_role.ecs.arn}"
}
