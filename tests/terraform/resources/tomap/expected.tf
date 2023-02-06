resource "aws_security_group" "sample" {
  name = "${var.this}-sample"

  tags = merge(var.tags, tomap({ "Name" = format("%s-sample", var.this) }), {
    test_tag = "test_value"
  })
}
