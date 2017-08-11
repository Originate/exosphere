resource "aws_security_group" "bastion" {
  name        = "${var.name}-bastion"
  vpc_id      = "${var.vpc_id}"
  description = "Bastion security group (only SSH inbound access is allowed)"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "${var.name}-bastion"
    Environment = "${var.env}"
  }

  lifecycle {
    create_before_destroy = true
  }
}
