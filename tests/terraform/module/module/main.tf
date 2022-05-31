module "complete_sg" {
  source              = "terraform-aws-modules/security-group/aws"
  name                = "my-sg-test"
  vpc_id              = "some-vpc-id"
  use_name_prefix     = true
  ingress_cidr_blocks = ["10.10.0.0/16"]
  ingress_rules       = ["https-443-tcp"]
  tags = {
    yor_trace = "6c510ebe-125a-41c0-aaf0-c19ccc56d5a0"
  }
}