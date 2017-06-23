resource "aws_internet_gateway" "public" {
  vpc_id = "${var.vpc_id}"

  tags {
    Name = "${var.env}-public"
  }
}

resource "aws_subnet" "public" {
  vpc_id            = "${var.vpc_id}"
  cidr_block        = "${cidrsubnet(var.vpc_cidr, 8, count.index + 1)}"
  availability_zone = "${element(var.availability_zones, count.index)}"
  count             = "${length(var.availability_zones)}"

  tags {
    Name = "${var.env}-public-${element(var.availability_zones, count.index)}"
  }

  depends_on = ["aws_internet_gateway.public"]

  lifecycle {
    create_before_destroy = true
  }

  map_public_ip_on_launch = true
}

resource "aws_route_table" "public" {
  vpc_id = "${var.vpc_id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.public.id}"
  }

  tags {
    Name = "${var.env}-public.${element(var.availability_zones, count.index)}"
  }
}

resource "aws_route_table_association" "public" {
  count          = "${length(var.availability_zones)}"
  subnet_id      = "${element(aws_subnet.public.*.id, count.index)}"
  route_table_id = "${aws_route_table.public.id}"
}

resource "aws_network_acl" "public" {
  vpc_id     = "${var.vpc_id}"
  subnet_ids = ["${aws_subnet.public.*.id}"]

  depends_on = ["aws_internet_gateway.public"]

  tags {
    Name = "${var.env}-public"
  }
}

resource "aws_network_acl_rule" "public-inbound" {
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = 100
  egress         = false
  protocol       = "all"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 0
  to_port        = 65535
}

resource "aws_network_acl_rule" "public-outbound" {
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = 100
  egress         = true
  protocol       = "all"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 0
  to_port        = 65535
}
