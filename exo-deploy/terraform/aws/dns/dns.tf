/**
 * The dns module creates a local route53 zone that serves
 * as a service discovery utility. For example a service
 * resource with the name `auth` and a dns module
 * with the name `stack.local`, the service address will be `auth.stack.local`.
 *
 * Usage:
 *
 *    module "dns" {
 *      source = "./dns"
 *      name   = "stack.local"
 *    }
 *
 */
resource "aws_route53_zone" "main" {
  name    = "${var.name}"
  vpc_id  = "${var.vpc_id}"
  comment = ""
}
