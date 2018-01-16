variable "{{password-secret-name}}" {}

module "{{db-name}}_rds_instance" {
  source = "github.com/Originate/exosphere.git//remote-dependency-templates//rds//modules//rds?ref={{terraformCommitHash}}"

  allocated_storage       = "{{allocated-storage}}"
  ecs_security_group      = "${data.terraform_remote_state.main_infrastructure.ecs_cluster_security_group}"
  bastion_security_group  = "${data.terraform_remote_state.main_infrastructure.bastion_security_group}"
  engine                  = "{{engine}}"
  engine_version          = "{{engine-version}}"
  env                     = "${var.env}"
  instance_class          = "{{instance-class}}"
  internal_hosted_zone_id = "${data.terraform_remote_state.main_infrastructure.internal_zone_id}"
  name                    = "{{db-name}}"
  username                = "{{username}}"
  password                = "${var.{{password-secret-name}}}"
  storage_type            = "{{storage-type}}"
  subnet_ids              = ["${data.terraform_remote_state.main_infrastructure.private_subnet_ids}"]
  vpc_id                  = "${data.terraform_remote_state.main_infrastructure.vpc_id}"
}
