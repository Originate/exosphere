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

//needs to be elsewhere (maybe cluster level)

resource "aws_alb" "alb" {
  name            = "${var.env}-${var.name}-lb"
  subnets         = ["${var.public_subnet_ids}"]
  security_groups = ["${aws_security_group.service_sg.id}"]
  internal        = false
}

resource "aws_alb_target_group" "target_group" {
  name     = "${var.env}-${var.name}"
  port     = 80
  protocol = "HTTP"
  vpc_id   = "${var.vpc_id}"

  health_check = {
    path = "/config.json"
    /* port = "${var.container_port}" */
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

output "security_groups" {
  value = ["${aws_security_group.service_sg.id}"]
}
