
module "{{name}}" {
  source = "./public-services"

  name = "{{name}}"
  cluster_id = "${module.public-cluster.cluster_id}"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.public_subnet_id}"
  public_port = {{publicPort}}
  service_type = "{{name}}"
}

resource "aws_route53_record" "{{name}}" {
  zone_id = "${var.hosted_zone_id}"
  name = "{{publicUrl}}"
  type = "A"

  alias {
  name = "${module.{{name}}.elb_dns_name}"
  zone_id = "${module.{{name}}.elb_zone_id}"
  evaluate_target_health = true
  }
}

output "{{name}}_url" { value = "${aws_route53_record.{{name}}.name}"}
