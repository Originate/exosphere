variable "{{passwordSecretName}}" {}

module "{{name}}_rds_instance" {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//dependencies//rds?ref={{terraformCommitHash}}"

  allocated_storage       = {{allocatedStorage}}
  ecs_security_group      = "${module.aws.ecs_cluster_security_group}"
  engine                  = "{{engine}}"
  engine_version          = "{{engineVersion}}"
  env                     = "production"
  instance_class          = "{{instanceClass}}"
  internal_hosted_zone_id = "${module.aws.internal_zone_id}"
  name                    = "{{name}}"
  username                = "{{username}}"
  password                = "${var.{{passwordSecretName}}}"
  storage_type            = "{{storageType}}"
  subnet_ids              = ["${module.aws.private_subnet_ids}"]
  vpc_id                  = "${module.aws.vpc_id}"
}
