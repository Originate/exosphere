resource "aws_security_group" "exocom_cluster" {
  name        = "exocom-ecs-cluster"
  vpc_id      = "${var.vpc_id}"
  description = "Allows traffic from and to the EC2 instances of the Exocom ECS cluster"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = -1
    security_groups = ["${var.security_groups}", "${aws_security_group.external_alb.id}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "ECS cluster (Exocom)"
    Environment = "${var.env}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

/* resource "aws_security_group" "external_alb" { */
/*   name        = "${format("exocom-%s-external-alb", var.env)}" */
/*   vpc_id      = "${var.vpc_id}" */
/*   description = "Allows external ALB traffic" */
/*  */
/*   ingress { */
/*     from_port   = 80 */
/*     to_port     = 80 */
/*     protocol    = "tcp" */
/*     cidr_blocks = ["0.0.0.0/0"] */
/*   } */
/*  */
/*   ingress { */
/*     from_port   = 443 */
/*     to_port     = 443 */
/*     protocol    = "tcp" */
/*     cidr_blocks = ["0.0.0.0/0"] */
/*   } */
/*  */
/*   egress { */
/*     from_port   = 0 */
/*     to_port     = 0 */
/*     protocol    = -1 */
/*     cidr_blocks = ["0.0.0.0/0"] */
/*   } */
/*  */
/*   lifecycle { */
/*     create_before_destroy = true */
/*   } */
/*  */
/*   tags { */
/*     Name        = "Exocom external ALB" */
/*     Environment = "${var.env}" */
/*   } */
/* } */
