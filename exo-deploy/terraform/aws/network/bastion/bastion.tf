variable "region" {
  description = "Region of the environment, for example, us-west-2"
}

variable "env" {
  description = "Name of the environment, used for naming and prefixing"
}

variable "vpc_id" {
  description = "ID of the VPC"
}

variable "availability_zones" {
  type        = "list"
  description = "List of availability zones to place subnets"
}

variable "public_subnet_ids" {
  type        = "list"
  description = "List of ID's of the public subnets"
}

variable "instance_type" {
  default     = "t2.micro"
  description = "Instance type of the bastion hosts"
}

variable "key_name" {
  description = "Name of the SSH key pair stored in AWS to authorize for the bastion hosts"
}

module "ami" {
  source        = "github.com/terraform-community-modules/tf_aws_ubuntu_ami/ebs"
  region        = "${var.region}"
  distribution  = "xenial"
  instance_type = "${var.instance_type}"
  storagetype   = "ebs-ssd"
}

resource "aws_security_group" "bastion" {
  name        = "${var.env}-bastion"
  vpc_id      = "${var.vpc_id}"
  description = "Bastion security group (only SSH inbound access is allowed)"

  ingress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "${var.env}-bastion"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_instance" "bastion" {
  ami           = "${module.ami.ami_id}"
  instance_type = "${var.instance_type}"

  vpc_security_group_ids = ["${aws_security_group.bastion.id}"]

  availability_zone           = "${element(var.availability_zones, count.index)}"
  subnet_id                   = "${element(var.public_subnet_ids, count.index)}"
  key_name                    = "${var.key_name}"
  associate_public_ip_address = true

  monitoring = true

  count = "${length(var.availability_zones)}"

  tags {
    Name = "${var.env}-bastion-${element(var.availability_zones, count.index)}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

output "security_group_id" {
  value = "${aws_security_group.bastion.id}"
}

output "private_ips" {
  value = ["${aws_instance.bastion.*.private_ip}"]
}

output "public_ips" {
  value = ["${aws_instance.bastion.*.public_ip}"]
}
