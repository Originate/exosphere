resource "aws_vpc_dhcp_options" "dns_resolver" {
  domain_name         = "${var.name}"
  domain_name_servers = ["${split(",", var.servers)}"]
}

resource "aws_vpc_dhcp_options_association" "dns_resolver" {
  vpc_id          = "${var.vpc_id}"
  dhcp_options_id = "${aws_vpc_dhcp_options.dns_resolver.id}"
}
