resource "aws_route53_zone" "zone" {
  name    = "${var.name}"
  vpc_id  = "${var.vpc_id}"
  comment = ""
}
