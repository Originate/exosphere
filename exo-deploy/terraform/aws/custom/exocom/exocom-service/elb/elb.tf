resource "aws_elb" "elb" {
  name            = "${var.env}-${var.name}-lb"
  /* instances = */
  internal        = false
  security_groups = ["${var.security_groups}"]
  subnets         = ["${var.subnet_ids}"]

  listener {
    instance_port     = "${var.host_port}"
    instance_protocol = "tcp"
    lb_port           = "80"
    lb_protocol       = "tcp"
  }

  /* health_check { */
  /*   healthy_threshold  = 5 */
  /*   unhealth_threshold = 3 */
  /*   target             = "HTTP:80/config.json" */
  /*   interval           = 30 */
  /* } */

}
