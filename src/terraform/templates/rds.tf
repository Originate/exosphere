variable "{{passwordEnvVar}}" {}

module "{{name}}_rds_instance" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//dependencies//rds?ref={{terraformCommitHash}}"

  allocated_storage       = {{allocatedStorage}}
  ecs_secuirty_group      = "${module.aws.ecs_cluster_security_group}"
  engine                  = "{{engine}}"
  engine_version          = "{{engineVersion}}"
  env                     = "production"
  instance_class          = "{{instanceClass}}"
  internal_hosted_zone_id = "${module.aws.internal_zone_id}"
  name                    = "{{name}}"
  username                = "{{username}}"
  password                = "${var.{{passwordEnvVar}}}"
  storage_type            = "{{storageType}}"
  subnet_ids              = ["${module.aws.private_subnet_ids}"]
}
