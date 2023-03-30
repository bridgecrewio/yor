locals {
  subnets = {
    "us-east-1" = {
      cidr_block = "10.10.10.10/24"
      tags       = {
        location = "us-east-1a"
      }
    }
  }
}

resource "aws_subnet" "eks_subnet" {
  for_each = local.subnets

  vpc_id                  = var.vpc_id
  cidr_block              = each.value.cidr_block
  availability_zone       = each.key
  map_public_ip_on_launch = true
  tags                    = each.value.tags
}