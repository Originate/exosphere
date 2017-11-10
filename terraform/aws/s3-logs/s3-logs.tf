data "aws_elb_service_account" "elb" {}

data "template_file" "policy" {
  template = "${file("${path.module}/files/policy.json")}"

  vars = {
    bucket        = "${var.name}-logs"
    principal_arn = "${data.aws_elb_service_account.elb.arn}"
  }
}

resource "aws_s3_bucket" "logs" {
  bucket        = "${var.name}-logs"
  policy        = "${data.template_file.policy.rendered}"
  force_destroy = true

  lifecycle_rule {
    id      = "logs-expiration"
    prefix  = ""
    enabled = "${var.logs_expiration_enabled}"

    expiration {
      days = "${var.logs_expiration_days}"
    }
  }

  tags {
    Name        = "${var.name}-logs"
    Environment = "${var.env}"
  }
}
