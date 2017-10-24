resource "aws_db_instance" "rds" {
  allocated_storage         = "${var.allocated_storage}"
  engine                    = "${var.engine}"
  engine_version            = "${var.engine_version}"
  db_subnet_group_name      = "${aws_db_subnet_group.rds_group.id}"
  instance_class            = "${var.instance_class}"
  final_snapshot_identifier = "${var.name}-final-snapshot"
  name                      = "${var.name}"
  username                  = "${var.username}"
  password                  = "${var.password}"
  storage_type              = "${var.storage_type}"
  vpc_security_group_ids    = ["${aws_security_group.rds.id}"]

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

resource "aws_route53_record" "rds" {
  zone_id = "${var.internal_hosted_zone_id}"
  name    = "${var.name}"
  type    = "CNAME"
  ttl     = "300"
  records = ["${aws_db_instance.rds.address}"]
}

resource "aws_security_group" "rds" {
  name        = "${var.name}-rds"
  vpc_id      = "${var.vpc_id}"
  description = "Allows traffic from ECS cluster"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = -1
    security_groups = ["${var.ecs_security_group}"]
  }

  ingress {
    from_port       = 5432
    to_port         = 0
    protocol        = -1
    security_groups = ["${var.bastion_security_group}"]
  }

  tags {
    Name        = "${var.name}-rds"
    Environment = "${var.env}"
  }

  lifecycle {
    create_before_destroy = true
  }
}
