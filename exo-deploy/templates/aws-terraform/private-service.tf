
module "{{name}}-service" {
  source = "../private-services"

  name = "exosphere-{{name}}-service"
  cluster_id = "${module.public-cluster.cluster_id}"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.public_subnet_id}" /* TODO: make private */
  service_type = "{{name}}"
}
