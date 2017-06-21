/* Get ECS optimized AMI id to use on the cluster */
data "aws_ami" "ecs_optimized_ami" {
  most_recent = true

  filter {
    name   = "owner-alias"
    values = ["amazon"]
  }

  filter {
    name   = "name"
    values = ["amzn-ami-2016.09.g-amazon-ecs-optimized"]
  }
}
