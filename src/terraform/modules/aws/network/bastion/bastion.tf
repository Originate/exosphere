module "ami" {
  source        = "github.com/terraform-community-modules/tf_aws_ubuntu_ami/ebs"
  region        = "${var.region}"
  distribution  = "xenial"
  instance_type = "${var.instance_type}"
  storagetype   = "ebs-ssd"
}

resource "aws_instance" "bastion" {
  ami           = "${module.ami.ami_id}"
  instance_type = "${var.instance_type}"

  vpc_security_group_ids = ["${aws_security_group.bastion.id}"]

  availability_zone           = "${element(var.availability_zones, count.index)}"
  subnet_id                   = "${element(var.subnet_ids, count.index)}"
  key_name                    = "${var.key_name}"
  associate_public_ip_address = true
  source_dest_check           = false
  monitoring                  = true

  count = "${length(var.availability_zones)}"

  tags {
    Name        = "${var.name}-bastion-${element(var.availability_zones, count.index)}"
    Environment = "${var.env}"
  }

  lifecycle {
    create_before_destroy = true
  }
}
