resource "aws_alb" "alb" {
  name            = "${var.env}-${var.name}-lb"
  subnets         = ["${var.subnet_ids}"]
  security_groups = ["${var.security_group}"]
  internal        = false

  access_logs {
    bucket = "${var.log_bucket}"
  }
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

resource "aws_route53_record" "internal" {
  zone_id = "${var.internal_zone_id}"
  name    = "${var.name}.${var.internal_dns_name}"
  type    = "A"

  alias {
    name                   = "${aws_alb.alb.dns_name}"
    zone_id                = "${aws_alb.alb.zone_id}"
    evaluate_target_health = false
  }
}
