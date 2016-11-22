variable "name" {}
variable "instance_type" {}
variable "security_groups" {}
variable "iam_instance_profile" {}
variable "subnet_id" {}

output "cluster_id" { value = "${aws_ecs_cluster.cluster.id}" }
output "exocom_ip"  { value = "${aws_instance.exocom.public_ip}" }


resource "aws_ecs_cluster" "cluster" {
  name = "exosphere-${var.name}-cluster"
}

resource "aws_instance" "exocom" {
  ami = "ami-56ed4936"
  subnet_id = "${var.subnet_id}"
  instance_type = "${var.instance_type}"
  iam_instance_profile = "${var.iam_instance_profile}"
  vpc_security_group_ids = ["${var.security_groups}"]
  user_data = "#!/bin/bash\necho ECS_CLUSTER=${aws_ecs_cluster.cluster.name} > /etc/ecs/ecs.config"
  key_name = "hugo"

  tags {
    Name = "exosphere-${var.name}-instance"
  }
}
