resource "aws_subnet" "private" {
  vpc_id            = "${var.vpc_id}"
  cidr_block        = "${cidrsubnet(var.vpc_cidr, 8, count.index + length(var.availability_zones) + 1)}"
  availability_zone = "${element(var.availability_zones, count.index)}"

  count = "${length(var.availability_zones)}"

  tags {
    Name = "${var.env}-private-${element(var.availability_zones, count.index)}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route_table" "private" {
  vpc_id = "${var.vpc_id}"
  count  = "${length(var.availability_zones)}"

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = "${element(aws_nat_gateway.nat.*.id, count.index)}"
  }

  tags {
    Name = "${var.env}-private-${element(var.availability_zones, count.index)}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route_table_association" "private" {
  count          = "${length(var.availability_zones)}"
  subnet_id      = "${element(aws_subnet.private.*.id, count.index)}"
  route_table_id = "${element(aws_route_table.private.*.id, count.index)}"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_network_acl" "private" {
  vpc_id     = "${var.vpc_id}"
  subnet_ids = ["${aws_subnet.private.*.id}"]

  tags {
    Name = "${var.env}-private"
  }
}

resource "aws_network_acl_rule" "private-inbound" {
  network_acl_id = "${aws_network_acl.private.id}"
  rule_number    = 100
  egress         = false
  protocol       = "all"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 0
  to_port        = 65535
}

resource "aws_network_acl_rule" "private-outbound" {
  network_acl_id = "${aws_network_acl.private.id}"
  rule_number    = 100
  egress         = true
  protocol       = "all"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 0
  to_port        = 65535
}

output "private_subnet_ids" {
  value = ["${aws_subnet.private.*.id}"]
}

output "private_subnet_cidrs" {
  value = ["${aws_subnet.private.*.cidr_block}"]
}
