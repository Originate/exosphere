variable "{{password-secret-name}}" {}

module "{{db-name}}_rds_instance" {
  source = "github.com/Originate/exosphere.git//remote-dependency-templates//rds//modules//rds?ref={{terraformCommitHash}}"

  allocated_storage       = "{{allocated-storage}}"
  ecs_security_group      = "${module.aws.ecs_cluster_security_group}"
  bastion_security_group  = "${module.aws.bastion_security_group}"
  engine                  = "{{engine}}"
  engine_version          = "{{engine-version}}"
  env                     = "${var.env}"
  instance_class          = "{{instance-class}}"
  internal_hosted_zone_id = "${module.aws.internal_zone_id}"
  name                    = "{{db-name}}"
  username                = "{{username}}"
  password                = "${var.{{password-secret-name}}}"
  storage_type            = "{{storage-type}}"
  subnet_ids              = ["${module.aws.private_subnet_ids}"]
  vpc_id                  = "${module.aws.vpc_id}"
}
