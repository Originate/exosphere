resource "aws_ecs_cluster" "exocom" {
  name = "exocom"

}

resource "aws_instance" "exocom" {
  ami                    = "${data.aws_ami.ecs_optimized.id}"
  iam_instance_profile   = "${aws_iam_instance_profile.exocom_ecs_instance.name}"
  instance_type          = "${var.instance_type}"
  key_name               = "${var.key_name}"
  security_groups        = ["${aws_security_group.exocom_cluster.id}"]
  subnet_id              = "${element(var.subnet_ids, length(var.subnet_ids))}"
  user_data              = "${data.template_cloudinit_config.cloud_config.rendered}"

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "${var.name}"
    Environment = "${var.env}"
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
