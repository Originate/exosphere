/* Variables */

variable "env" {
  description = "Environment tag, e.g prod"
}

variable "cidr" {
  description = "The cidr block for the VPC"
  default     = "10.0.0.0/16"
}

variable "name" {
  default = "vpc"
}

/* Output */

output "id" {
  value = "${aws_vpc.vpc.id}"
}

output "cidr" {
  value = "${aws_vpc.vpc.cidr_block}"
}
