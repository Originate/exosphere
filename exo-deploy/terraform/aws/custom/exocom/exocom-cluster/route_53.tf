resource "aws_route53_record" "exocom" {
  zone_id = "${var.internal_hosted_zone_id}"
  name    = "${var.name}.${var.internal_dns_name}"
  type    = "A"
  ttl     = "300"
  records = ["${aws_instance.exocom.private_ip}"]
}
