data "aws_ami" "ecs_optimized" {
  most_recent = true

  filter {
    name   = "owner-alias"
    values = ["amazon"]
  }

  filter {
    name   = "name"
    values = ["amzn-ami-2017.03.d-amazon-ecs-optimized"]
  }
}
