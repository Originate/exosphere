resource "aws_security_group" "exocom_cluster" {
  name        = "exocom-ecs-cluster"
  vpc_id      = "${var.vpc_id}"
  description = "Allows traffic from and to the EC2 instances of the Exocom ECS cluster"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = -1
    security_groups = ["${var.security_groups}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "ECS cluster (Exocom)"
    Environment = "${var.env}"
  }

  lifecycle {
    create_before_destroy = true
  }
}
