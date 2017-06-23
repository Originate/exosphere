resource "aws_ecs_cluster" "exocom" {
  name = "exocom"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_launch_configuration" "exocom" {
  name_prefix = "exocom-"

  image_id             = "${data.aws_ami.ecs_optimized_ami.id}"
  instance_type        = "${var.instance_type}"
  ebs_optimized        = "${var.ebs_optimized}"
  iam_instance_profile = "${aws_iam_instance_profile.exocom_ecs_instance.arn}"
  security_groups      = ["${aws_security_group.exocom_cluster.id}"]
  key_name             = "${var.key_name}"

  user_data = "#!/bin/bash\necho ECS_CLUSTER=${aws_ecs_cluster.exocom.name} > /etc/ecs/ecs.config"
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

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/files/cloud-config.yml.tpl")}"

  vars {
    environment      = "${var.env}"
    name             = "exocom"
    region           = "${var.region}"
    docker_auth_type = "${var.docker_auth_type}"
    docker_auth_data = "${var.docker_auth_data}"
  }
}

data "template_cloudinit_config" "cloud_config" {
  gzip          = false
  base64_encode = false

  part {
    content_type = "text/cloud-config"
    content      = "${data.template_file.ecs_cloud_config.rendered}"
  }

  part {
    content_type = "${var.extra_cloud_config_type}"
    content      = "${var.extra_cloud_config_content}"
  }
}
