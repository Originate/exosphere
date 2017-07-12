resource "aws_ecs_cluster" "exocom" {
  name = "exocom"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_instance" "exocom" {
  ami                         = "${data.aws_ami.ecs_optimized_ami.id}"
  associate_public_ip_address = true
  iam_instance_profile        = "${aws_iam_instance_profile.exocom_ecs_instance.name}"
  instance_type               = "${var.instance_type}"
  key_name                    = "${var.key_name}"
  security_groups             = ["${aws_security_group.exocom_cluster.id}"]
  subnet_id                   = "${element(var.subnet_ids, length(var.subnet_ids))}"
  vpc_security_group_ids      = ["${var.security_groups}"]

  user_data = "#!/bin/bash\necho ECS_CLUSTER=${aws_ecs_cluster.exocom.name} > /etc/ecs/ecs.config"

  tags {
    Name = "exocom-instance"
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
