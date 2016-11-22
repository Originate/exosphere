output "space_tweet_url" { value = "${aws_route53_record.web.name}"}

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

module "exocom-cluster" {
  source = "./exocom-cluster"

  name = "exocom"
  iam_instance_profile = "ecsInstanceRole"
  instance_type = "t2.micro"
  security_groups = "${aws_security_group.public.id}"
  subnet_id = "${module.vpc.public_subnet_id}"
}

module "exocom-service" {
  source = "./exocom-service"

  name = "exosphere-exocom"
  cluster_id = "${module.exocom-cluster.cluster_id}"
}

module "public-cluster" {
  source = "./public-cluster"

  name = "space-tweet-public-cluster"
  iam_instance_profile = "ecsInstanceRole"
  instance_type = "t2.medium"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.public_subnet_id}"
}

resource "aws_route53_record" "exocom" {
   zone_id = "ZLQY1E06TWU1R"
   name = "exocomm.us-west-1a.aws.megabyt.es."
   type = "A"
   ttl = "300"
   records = ["${module.exocom-cluster.exocom_ip}"]
}

module "users-service" {
  source = "./private-services"

  name = "exosphere-users-service"
  cluster_id = "${module.public-cluster.cluster_id}"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.public_subnet_id}"
  service_type = "users-service"
}

module "tweets-service" {
  source = "./private-services"

  name = "exosphere-tweets-service"
  cluster_id = "${module.public-cluster.cluster_id}"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.public_subnet_id}"
  service_type = "tweets-service"
}

module "web-service" {
  source = "./public-services"

  name = "exosphere-web-service"
  cluster_id = "${module.public-cluster.cluster_id}"
  security_groups = "${aws_security_group.public.id}"
  subnet_ids = "${module.vpc.public_subnet_id}"
  service_type = "web-service"
}

resource "aws_route53_record" "web" {
  zone_id = "ZLQY1E06TWU1R"
  name = "megabyt.es."
  type = "A"

  alias {
  name = "${module.web-service.web_elb_dns_name}"
  zone_id = "${module.web-service.web_elb_zone_id}"
  evaluate_target_health = true
  }
}
