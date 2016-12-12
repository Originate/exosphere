module "vpc" {
  source = "./vpc"
}


resource "aws_security_group" "public" {
  name = "ECS-public-sg"
  vpc_id = "${module.vpc.vpc_id}"

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
      from_port = 0
      to_port = 0
      protocol = "-1"
      cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_route53_zone" "hosted_zone" {
  name = "{{domainName}}"
}


data "aws_ami" "ecs_ami" {
  most_recent = true
  filter {
    name = "owner-alias"
    values = ["amazon"]
  }
  filter {
    name = "name"
    values = ["amzn-ami-*-amazon-ecs-optimized"]
  }
}


module "exocom-cluster" {
  source = "./exocom/cluster"

  name = "exocom"
  iam_instance_profile = "ecsInstanceRole"
  instance_type = "t2.micro"
  security_groups = "${aws_security_group.public.id}"
  subnet_id = "${module.vpc.public_subnet_id}"
  ami_id = "${data.aws_ami.ecs_ami.id}"
}


module "public-cluster" {
  source = "./public-services/cluster"

  name = "{{appName}}-public-cluster"
  iam_instance_profile = "ecsInstanceRole"
  instance_type = "{{publicClusterSize}}"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.public_subnet_id}"
  ami_id = "${data.aws_ami.ecs_ami.id}"
}

module "private-cluster" {
  source = "./private-services/cluster"

  name = "{{appName}}-private-cluster"
  iam_instance_profile = "ecsInstanceRole"
  instance_type = "{{privateClusterSize}}"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.private_subnet_id}"
  ami_id = "${data.aws_ami.ecs_ami.id}"
}
