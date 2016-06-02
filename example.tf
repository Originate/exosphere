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


module "public-cluster" {
  source = "./cluster"

  name = "space-tweet-public-cluster"
  iam_instance_profile = "ecsInstanceRole"
  instance_type = "t2.micro"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.public_subnet_ids}"
}


module "private-cluster" {
  source = "./cluster"

  name = "space-tweet-private-cluster"
  iam_instance_profile = "ecsInstanceRole"
  instance_type = "t2.micro"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.private_subnet_ids}"
}


module "service" {
  source = "./service"

  name = "space-tweet-user-service"
  cluster_id = "${module.public-cluster.cluster_id}"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.public_subnet_ids}"
}
