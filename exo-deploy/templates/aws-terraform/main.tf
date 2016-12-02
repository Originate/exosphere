resource "aws_route53_zone" "hosted_zone" {
  name = "{{domainName}}"
}

module "main-infrastructure" {
  source = "./main-infrastructure"

  hosted_zone_id = "${aws_route53_zone.hosted_zone.zone_id}"
}
