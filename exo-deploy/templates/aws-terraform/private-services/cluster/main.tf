variable "name" {}
variable "instance_type" {}
variable "security_groups" {}
variable "iam_instance_profile" {}
variable "subnet_ids" {}
variable "ami_id" {}

output "cluster_id" { value = "${aws_ecs_cluster.cluster.id}" }


resource "aws_ecs_cluster" "cluster" {
  name = "exosphere-${var.name}-cluster"
}


/* Template used by auto-scaling group to launch EC2 instances. */
resource "aws_launch_configuration" "cluster" {
  name                 = "exosphere-${var.name}-launch-config"
  iam_instance_profile = "${var.iam_instance_profile}"
  image_id             = "${var.ami_id}"
  instance_type        = "${var.instance_type}"
  security_groups      = ["${var.security_groups}"]
  /* This command in user_data ensures the machine that's spun up with this
  launch configuration is associated with the appropriate ECS cluster. */
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
