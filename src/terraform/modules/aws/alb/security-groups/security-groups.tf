resource "aws_security_group" "internal_alb" {
  name        = "${var.name}-internal-alb"
  vpc_id      = "${var.vpc_id}"
  description = "Allows internal VPC traffic through the ALB"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "${var.name}-internal-alb"
    Environment = "${var.env}"
  }
}

resource "aws_security_group" "external_alb" {
  name        = "${var.name}-external-alb"
  vpc_id      = "${var.vpc_id}"
  description = "Allows all external traffic through the ALB"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "${var.name}-external-alb"
    Environment = "${var.env}"
  }
}
