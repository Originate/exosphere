resource "aws_security_group" "service_sg" {
  name   = "${var.env}-${var.name}-sg"
  vpc_id = "${var.vpc_id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_alb" "alb" {
  name            = "${var.env}-${var.name}-lb"
  subnets         = ["${var.subnet_ids}"]
  security_groups = ["${aws_security_group.service_sg.id}"]
  internal        = false
}

resource "aws_alb_target_group" "target_group" {
  name     = "${var.env}-${var.name}"
  port     = 80
  protocol = "HTTP"
  vpc_id   = "${var.vpc_id}"

  health_check = {
    path = "${var.health_check_endpoint}"
  }
}

resource "aws_alb_listener" "listener" {
  load_balancer_arn = "${aws_alb.alb.arn}"
  port              = "80"
  protocol          = "HTTP"

  default_action {
    target_group_arn = "${aws_alb_target_group.target_group.arn}"
    type             = "forward"
  }
}

output "url" {
  value = "${aws_alb.alb.dns_name}"
}

output "target_group_id" {
  value = "${aws_alb_target_group.target_group.id}"
}

output "security_groups" {
  value = ["${aws_security_group.service_sg.id}"]
}
