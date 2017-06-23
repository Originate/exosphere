resource "aws_vpc" "vpc" {
  cidr_block           = "${var.cidr}"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags {
    Name = "${var.env}-${var.name}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

output "id" {
  value = "${aws_vpc.vpc.id}"
}

output "cidr" {
  value = "${aws_vpc.vpc.cidr_block}"
}
