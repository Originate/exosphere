data "template_file" "policy" {
  template = "${file("${path.module}/policy.json")}"

  vars = {
    bucket     = "${var.name}-${var.env}-logs"
    account_id = "${var.account_id}"
  }
}

resource "aws_s3_bucket" "logs" {
  bucket = "${var.name}-${var.env}-logs"

  lifecycle_rule {
    id = "logs-expiration"
    prefix = ""
    enabled = "${var.logs_expiration_enabled}"

    expiration {
      days = "${var.logs_expiration_days}"
    }
  }

  tags {
    Name        = "${var.name}-${var.env}-logs"
    Environment = "${var.env}"
  }

  policy = "${data.template_file.policy.rendered}"
}
