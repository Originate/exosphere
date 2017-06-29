data "aws_elb_service_account" "elb" {}

data "template_file" "policy" {
  template = "${file("${path.module}/policy.json")}"

  vars = {
    bucket        = "${var.name}-${var.env}-logs"
    principal_arn = "${data.aws_elb_service_account.elb.arn}"
  }
}

resource "aws_s3_bucket" "logs" {
  bucket        = "${var.name}-${var.env}-logs"
  force_destroy = true
  policy        = "${data.template_file.policy.rendered}"

  lifecycle_rule {
    id      = "logs-expiration"
    prefix  = ""
    enabled = "${var.logs_expiration_enabled}"

    expiration {
      days = "${var.logs_expiration_days}"
    }
  }

  tags {
    Name        = "${var.name}-${var.env}-logs"
    Environment = "${var.env}"
  }
}
