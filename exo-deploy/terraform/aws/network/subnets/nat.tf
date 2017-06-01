resource "aws_eip" "nat" {
  vpc = true

  count = "${length(var.availability_zones)}"

  depends_on = ["aws_internet_gateway.public"]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_nat_gateway" "nat" {
  allocation_id = "${element(aws_eip.nat.*.id, count.index)}"
  subnet_id     = "${element(aws_subnet.public.*.id, count.index)}"

  count = "${length(var.availability_zones)}"

  depends_on = ["aws_internet_gateway.public"]

  lifecycle {
    create_before_destroy = true
  }
}

output "nat_gateway_ids" {
  value = ["${aws_nat_gateway.nat.*.id}"]
}
