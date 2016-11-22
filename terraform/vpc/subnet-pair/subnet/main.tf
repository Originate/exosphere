variable "availability_zone" {}
variable "cidr_block" {}
variable "name" {}
variable "vpc_id" {}

output "route_table_id" { value = "${aws_route_table.subnet-route-table.id}" }
output "subnet_id"      { value = "${aws_subnet.subnet.id}" }


resource "aws_subnet" "subnet" {
  vpc_id = "${var.vpc_id}"
  cidr_block = "${var.cidr_block}"
  availability_zone = "${var.availability_zone}"
  map_public_ip_on_launch = true

  tags {
    Name = "${var.name}"
  }
}


resource "aws_route_table" "subnet-route-table" {
  vpc_id = "${var.vpc_id}"

  tags {
    Name = "${var.name}-route-table"
  }
}


resource "aws_route_table_association" "route-table-association" {
  subnet_id = "${aws_subnet.subnet.id}"
  route_table_id = "${aws_route_table.subnet-route-table.id}"
}
