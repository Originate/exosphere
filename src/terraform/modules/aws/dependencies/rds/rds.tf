resource "aws_db_instance" "rds" {
  allocated_storage      = "${var.allocated_storage}"
  engine                 = "${var.engine}"
  engine_version         = "${var.engine_version}"
  db_subnet_group_name   = "${aws_db_subnet_group.rds_group.id}"
  instance_class         = "${var.instance_class}"
  name                   = "${var.name}"
  username               = "${var.username}"
  password               = "${var.password}"
  storage_type           = "${var.storage_type}"
  vpc_security_group_ids = "${var.vpc_security_group_ids}"

  tags {
    Name        = "${var.name}"
    Environment = "${var.env}"
  }
}
resource "aws_db_subnet_group" "rds_group" {
  name       = "${var.name}"
  subnet_ids = ["${var.subnet_ids}"]

  tags {
    Name        = "${var.name}"
    Environment = "${var.env}"
  }
}
