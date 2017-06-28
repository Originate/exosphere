resource "aws_route53_record" "exocom" {
  zone_id = "${var.hosted_zone_id}"
  name = "${format("exocom.%s", var.domain_name)}"
  type = "A"
  ttl = "300"
  records = ["${aws_instance.exocom.private_ip}"]
}
