variable "availability_zones" {
  description = "Availability zones to use for subnets. Two subnets will be created per availability zone"
  type        = "list"
}

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "vpc_cidr" {
  description = "The cidr block of the VPC"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

output "nat_gateway_ids" {
  value = ["${aws_nat_gateway.nat.*.id}"]
}

output "private_subnet_ids" {
  value = ["${aws_subnet.private.*.id}"]
}

output "private_subnet_cidrs" {
  value = ["${aws_subnet.private.*.cidr_block}"]
}

output "internet_gateway_id" {
  value = "${aws_internet_gateway.public.id}"
}

output "public_subnet_ids" {
  value = ["${aws_subnet.public.*.id}"]
}

output "public_subnet_cidrs" {
  value = ["${aws_subnet.public.*.cidr_block}"]
}
