output "vpc_id" { value = "${aws_vpc.exosphere.id}" }
output "public_subnet_id" {
  value = "${module.subnet-pair-2a.public_subnet_id}"
}
output "private_subnet_id" {
  value = "${module.subnet-pair-2a.private_subnet_id}"
}


resource "aws_vpc" "exosphere" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "exosphere-vpc"
  }
}


resource "aws_internet_gateway" "exosphere-igw" {
  vpc_id = "${aws_vpc.exosphere.id}"

  tags {
    Name = "exosphere-igw"
  }
}


module "subnet-pair-2a" {
  source = "./subnet-pair"
  availability_zone = "us-west-2a"
  gateway_id = "${aws_internet_gateway.exosphere-igw.id}"
  private_cidr_block = "10.0.32.0/19"
  public_cidr_block = "10.0.0.0/19"
  vpc_id = "${aws_vpc.exosphere.id}"
}
