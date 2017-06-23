resource "aws_security_group" "internal_alb" {
  name        = "${format("%s-%s-internal-alb", var.name, var.env)}"
  vpc_id      = "${var.vpc_id}"
  description = "Allows internal ALB traffic"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    /* cidr_blocks = ["${var.cidr}"] */
    cidr_blocks = ["0.0.0.0/0"]
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
    Name        = "${format("%s internal ALB", var.name)}"
    Environment = "${var.env}"
  }
}

resource "aws_security_group" "external_alb" {
  name        = "${format("%s-%s-external-alb", var.name, var.env)}"
  vpc_id      = "${var.vpc_id}"
  description = "Allows external ALB traffic"

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
    Name        = "${format("%s external ALB", var.name)}"
    Environment = "${var.env}"
  }
}
