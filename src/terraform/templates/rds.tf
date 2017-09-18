module "rds_instance"' {
  source = "git@github.com:Originate/exosphere.git//src//terraform//modules//aws//dependencies//exocom//exocom-cluster?ref={{terraformCommitHash}}"

  allocated_storage = "{{allocatedStorage}}"
  engine            = "{{engine}}"
  engine_version    = "{{engineVersion}}"
  env               = "production"
  instance_class    = "{{instanceClass}}"
  name              = "{{name}}"
  username          = "{{username}}"
  password          = "{{password}}"
  storage_type      = "{{storageType}}"
  subnet_ids        = ["${module.aws.private_subnet_ids}"]
}
