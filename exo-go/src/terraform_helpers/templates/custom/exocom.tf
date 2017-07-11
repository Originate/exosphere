module "exocom_cluster" {
  source                  = "./aws/custom/exocom/exocom-cluster"

  name                    = "exocom"

  availability_zones      = "${module.aws.availability_zones}"
  env                     = "production"
  internal_dns_name       = "${module.aws.internal_dns_name}"
  internal_hosted_zone_id = "${module.aws.internal_hosted_zone_id}"
  instance_type           = "t2.micro"
  key_name                = "${var.key_name}"
  region                  = "${var.region}"
  security_groups         = ["${module.aws.bastion_security_group_id}", "${module.aws.cluster_security_group}", "${module.aws.external_alb_security_group}"]
  subnet_ids              = "${module.aws.private_subnet_ids}"
  vpc_id                  = "${module.aws.vpc_id}"
}

module "exocom_service" {
  source                      = "./aws/custom/exocom/exocom-service"

  cluster_id                  = "${module.exocom_cluster.cluster_id}"
  command                     = ["bin/exocom"]
  container_port              = "3100"
  cpu_units                   = "128"
  docker_image                = "${var.account_id}.dkr.ecr.${var.region}.amazonaws.com/exocom:latest"
  env                         = "production"
  environment_variables       = {
    ROLE = "exocom"
    SERVICE_ROUTES = <<EOF
{{serviceRoutes}}
EOF
  }
  memory_reservation          = "128"
  name                        = "exocom"
  region                      = "${var.region}"
}
