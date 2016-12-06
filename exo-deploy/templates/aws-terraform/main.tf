resource "aws_route53_zone" "hosted_zone" {
  name = "{{domainName}}"
}

module "main-infrastructure" {
  source = "./main-infrastructure"

  hosted_zone_id = "${aws_route53_zone.hosted_zone.zone_id}"
}

output "hosted_zone_id" { value = "${aws_route53_zone.hosted_zone.zone_id}" }

#TODO: extract main-infrastructure.tf outputs here
