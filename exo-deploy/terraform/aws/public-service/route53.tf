resource "aws_route53_record" "public_url" {
  zone_id = "${var.external_hosted_zone_id}"
  name    = "${var.external_dns_name}"
  type    = "A"

  alias {
    name                   = "${module.external_alb.dns_name}"
    zone_id                = "${module.external_alb.zone_id}"
    evaluate_target_health = true
  }
}
