module "complete_sg" {
  source              = "terraform-aws-modules/security-group/aws"
  name                = "my-sg-test"
  vpc_id              = "some-vpc-id"
  use_name_prefix     = true
  ingress_cidr_blocks = ["10.10.0.0/16"]
  ingress_rules       = ["https-443-tcp"]

  tags = {
    Name      = "test-sg"
    yor_trace = "c8f514ca-29cd-479d-bee0-faead1d8ac2e"
  }
}