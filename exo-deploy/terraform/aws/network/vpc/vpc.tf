variable "env" {
  description = "Environment tag, e.g prod"
}

variable "name" {
  default = "vpc"
}

variable "cidr" {
  description = "The cidr block for the VPC"
  default     = "10.0.0.0/16"
}

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
