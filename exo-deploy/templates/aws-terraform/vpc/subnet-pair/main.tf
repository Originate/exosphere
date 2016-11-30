variable "availability_zone" {}
variable "gateway_id" {}
variable "private_cidr_block" {}
variable "public_cidr_block" {}
variable "vpc_id" {}

output "public_subnet_id" { value = "${module.public-subnet.subnet_id}" }
output "private_subnet_id" { value = "${module.private-subnet.subnet_id}" }


module "public-subnet" {
  source = "./subnet"

  name = "exosphere-${var.availability_zone}-public-subnet"
  availability_zone = "${var.availability_zone}"
  cidr_block = "${var.public_cidr_block}"
  vpc_id = "${var.vpc_id}"
}


module "private-subnet" {
  source = "./subnet"

  name = "exosphere-${var.availability_zone}-private-subnet"
  availability_zone = "${var.availability_zone}"
  cidr_block = "${var.private_cidr_block}"
  vpc_id = "${var.vpc_id}"
}


resource "aws_route" "public-subnet-igw-route" {
  destination_cidr_block = "0.0.0.0/0"
  gateway_id = "${var.gateway_id}"
  route_table_id = "${module.public-subnet.route_table_id}"
}
