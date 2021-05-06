resource "aws_security_group" "cluster" {
  name_prefix = "acme"
  description = "EKS cluster security group"
  vpc_id      = "vpc-123456"
  tags = merge(
  var.tags,
  {
    "Name" = "${var.env}-eks_cluster_sg"
  },
  )
}

resource "aws_vpc" "vpc_tags_one_line" {
  cidr_block = ""
  tags = { "Name" = "tag-for-s3", "Environment" = "prod" }
}

resource "aws_instance" "no_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"
}

resource "aws_instance" "simple_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"

  tags = {
    Name = "my-instance"
  }
}

resource "aws_instance" "rendered_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"

  tags = var.tags
}

resource "aws_instance" "merge_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"

  tags = merge(var.tags,
  {
    Name = "merged-tags-instance",
    Env  = var.env
  })
}

variable "tags" {
  default = {}
  type = map(string)
}

variable "env" {
  default = "dev"
  type = string
}