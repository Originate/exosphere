variable "name" {}
variable "instance_type" {}
variable "security_groups" {}
variable "iam_instance_profile" {}
variable "subnet_ids" {}

output "cluster_id" { value = "${aws_ecs_cluster.cluster.id}" }


resource "aws_ecs_cluster" "cluster" {
  name = "exosphere-${var.name}-cluster"
}


resource "aws_launch_configuration" "cluster" {
  name                 = "exosphere-${var.name}-launch-config"
  iam_instance_profile = "${var.iam_instance_profile}"
  image_id             = "ami-56ed4936"
  instance_type        = "${var.instance_type}"
  security_groups      = ["${var.security_groups}"]
  user_data            = "#!/bin/bash\necho ECS_CLUSTER=${aws_ecs_cluster.cluster.name} > /etc/ecs/ecs.config"
}


resource "aws_autoscaling_group" "cluster" {
  name = "exosphere-${var.name}-autoscaling-group"
  min_size = "${length(split(",", var.subnet_ids))}"
  max_size = "${length(split(",", var.subnet_ids))}"
  launch_configuration = "${aws_launch_configuration.cluster.name}"
  vpc_zone_identifier = ["${split(",", var.subnet_ids)}"]

  tag {
    key = "Name"
    value = "exosphere-${var.name}-instance"
    propagate_at_launch = true
  }
}
