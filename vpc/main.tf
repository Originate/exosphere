output "vpc_id" { value = "${aws_vpc.exosphere.id}" }
output "public_subnet_ids" {
  value = "${module.subnet-pair-2a.public_subnet_id},${module.subnet-pair-2b.public_subnet_id},${module.subnet-pair-2c.public_subnet_id}"
}
output "private_subnet_ids" {
  value = "${module.subnet-pair-2a.private_subnet_id},${module.subnet-pair-2b.private_subnet_id},${module.subnet-pair-2c.private_subnet_id}"
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


module "subnet-pair-2b" {
  source = "./subnet-pair"
  availability_zone = "us-west-2b"
  gateway_id = "${aws_internet_gateway.exosphere-igw.id}"
  public_cidr_block = "10.0.64.0/19"
  private_cidr_block = "10.0.96.0/19"
  vpc_id = "${aws_vpc.exosphere.id}"
}


module "subnet-pair-2c" {
  source = "./subnet-pair"
  availability_zone = "us-west-2c"
  gateway_id = "${aws_internet_gateway.exosphere-igw.id}"
  public_cidr_block = "10.0.128.0/19"
  private_cidr_block = "10.0.160.0/19"
  vpc_id = "${aws_vpc.exosphere.id}"
}
