resource "aws_vpc" "tagged_vpc" {
  cidr_block = "10.0.0.0/16"
  tags = merge(var.tags, {
    yor_trace = "3b06a351-83d4-4922-90b3-a4a35bab9bb7"
  })
}

resource "aws_s3_bucket" "my-bucket" {
  bucket = "my-bucket"
  tags = {
    yor_trace = "99205876-f974-4075-b032-12fac6bfa938"
  }
}

module "sg" {
  source = "terraform-aws-modules/vpc/aws"
  tags = {
    yor_trace = "bb1e4a74-94f5-4226-bf07-127ba92c2b89"
  }
}