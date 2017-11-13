resource "aws_route53_record" "internal" {
  zone_id = "${var.internal_zone_id}"
  name    = "${var.internal_dns_name}"
  type    = "A"

  alias {
    zone_id                = "${aws_alb.alb.zone_id}"
    name                   = "${aws_alb.alb.dns_name}"
    evaluate_target_health = false
  }
}
