resource "aws_security_group" "cluster" {
  name        = "${var.name}-ecs-cluster"
  vpc_id      = "${var.vpc_id}"
  description = "Allows traffic from and to the EC2 instances of the ${var.name} ECS cluster"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = -1
    security_groups = ["${var.security_groups}"] //allow acccess from alb
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "ECS cluster (${var.name})"
    Environment = "${var.env}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_ecs_cluster" "main" {
  name = "${var.name}"

  lifecycle {
    create_before_destroy = true
  }
}

/* data "template_file" "ecs_cloud_config" { */
/*   template = "${file("${path.module}/files/cloud-config.yml.tpl")}" */
/*  */
/*   vars { */
/*     environment      = "${var.env}" */
/*     name             = "${var.name}" */
/*     region           = "${var.region}" */
/*   } */
/* } */

/* data "template_cloudinit_config" "cloud_config" { */
/*   gzip          = false */
/*   base64_encode = false */
/*  */
/*   part { */
/*     content_type = "text/cloud-config" */
/*     content      = "${data.template_file.ecs_cloud_config.rendered}" */
/*   } */
/*  */
/*   part { */
/*     content_type = "${var.extra_cloud_config_type}" */
/*     content      = "${var.extra_cloud_config_content}" */
/*   } */
/* } */

module "iam" {
  source = "./iam"

  env    = "${var.env}"
}

resource "aws_launch_configuration" "main" {
  name_prefix = "${format("%s-", var.name)}"

  image_id             = "${var.image_id}"
  instance_type        = "${var.instance_type}"
  ebs_optimized        = false
  iam_instance_profile = "${module.iam.iam_instance_profile}"
  security_groups      = ["${aws_security_group.cluster.id}"]
  key_name             = "${var.key_name}"

  user_data = "#!/bin/bash\necho ECS_CLUSTER=${aws_ecs_cluster.main.name} > /etc/ecs/ecs.config"
  associate_public_ip_address = false

  # root
  root_block_device {
    volume_type = "gp2"
    volume_size = "${var.root_volume_size}"
  }

  # docker
  ebs_block_device {
    device_name = "/dev/xvdcz"
    volume_type = "gp2"
    volume_size = "${var.docker_volume_size}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "main" {
  name = "${var.name}"

  availability_zones   = ["${var.availability_zones}"]
  vpc_zone_identifier  = ["${var.subnet_ids}"]
  launch_configuration = "${aws_launch_configuration.main.id}"
  min_size             = "${var.min_size}"
  max_size             = "${var.max_size}"
  desired_capacity     = "${var.desired_capacity}"
  termination_policies = ["OldestLaunchConfiguration", "Default"]

  tag {
    key                 = "Name"
    value               = "${var.name}"
    propagate_at_launch = true
  }

  tag {
    key                 = "Cluster"
    value               = "${var.name}"
    propagate_at_launch = true
  }

  tag {
    key                 = "Environment"
    value               = "${var.env}"
    propagate_at_launch = true
  }

  lifecycle {
    create_before_destroy = true
  }
}

module "autoscaling_policy" {
  source = "./autoscaling_policy"

  name                   = "${var.name}"
  autoscaling_group_name = "${aws_autoscaling_group.main.name}"
  cluster_name           = "${aws_ecs_cluster.main.name}"
}

output "id" {
  value = "${aws_ecs_cluster.main.id}"
}

// The cluster security group ID.
output "security_group_id" {
  value = "${aws_security_group.cluster.id}"
}
