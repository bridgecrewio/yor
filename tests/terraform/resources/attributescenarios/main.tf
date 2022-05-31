resource "aws_security_group" "cluster" {
  name_prefix = "acme"
  description = "EKS cluster security group"
  vpc_id      = "vpc-123456"
  tags = merge(
    var.tags,
    {
      "Name" = "${var.env}-eks_cluster_sg"
    },
    {
      yor_trace = "912f3703-2158-48c1-be6d-e1c5b0c60eab"
  })
}

resource "aws_vpc" "vpc_tags_one_line" {
  cidr_block = ""
  tags = { "Name" = "tag-for-s3", "Environment" = "prod"
    yor_trace = "da94701e-a1b9-414c-b4b2-57573e6ba92d"
  }
}

resource "aws_instance" "no_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"
  tags = {
    yor_trace = "091e247f-4a6f-45de-85bb-015096075450"
  }
}

resource "aws_instance" "simple_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"

  tags = {
    Name      = "my-instance"
    yor_trace = "13bbfe58-2d72-43e5-84b0-75e63ce412eb"
  }
}

resource "aws_instance" "rendered_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"

  tags = merge(var.tags, {
    yor_trace = "ad8ff30e-b4b3-485c-bd39-1f9bc0a0d965"
  })
}

resource "aws_instance" "merge_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"

  tags = merge(var.tags,
    {
      Name = "merged-tags-instance",
      Env  = var.env
      }, {
      yor_trace = "292ebec9-4149-4651-90ed-ac39f229dd3b"
  })
}

variable "tags" {
  default = {}
  type    = map(string)
}

variable "env" {
  default = "dev"
  type    = string
}