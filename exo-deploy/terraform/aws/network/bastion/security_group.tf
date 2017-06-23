resource "aws_security_group" "bastion" {
  name        = "${var.env}-bastion"
  vpc_id      = "${var.vpc_id}"
  description = "Bastion security group (only SSH inbound access is allowed)"

  ingress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "${var.env}-bastion"
  }

  lifecycle {
    create_before_destroy = true
  }
}
