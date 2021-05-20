module "complete_sg" {
  source              = "terraform-aws-modules/security-group/aws"
  name                = "my-sg-test"
  vpc_id              = "some-vpc-id"
  use_name_prefix     = true
  ingress_cidr_blocks = ["10.10.0.0/16"]
  ingress_rules       = ["https-443-tcp"]

  tags = {
    Name = "test-sg"
  }
   = {
    Name      = "test-sg"
    yor_trace = "955a885c-9292-4288-bebb-ccbf5e613ef4"
  }
}