resource "aws_alb" "alb" {
  name            = "${var.env}-${var.name}-lb"
  subnets         = ["${var.subnet_ids}"]
  security_groups = ["${var.security_groups}"]
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
