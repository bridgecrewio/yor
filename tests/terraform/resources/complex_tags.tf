resource "aws_s3_bucket" "bucket_var_tags" {
  tags = var.tags
}

variable "tags" {
  default = {
    "Name": "tag-for-s3"
    "Environment":"prod"
  }
}

resource "aws_vpc" "vpc_tags_one_line" {
  cidr_block = ""
  tags = {"Name": "tag-for-s3", "Environment":"prod"}
}

resource "aws_alb" "alb_with_merged_tags" {
  tags = merge({"Name": "tag-for-s3", "Environment":"prod"},
  {"yor_trace": "4329587194",
    "git_org": "bana"})
}

resource "aws_instance" "many_instance_tags" {
  ami = ""
  instance_type = ""
  tags =  merge({"Name": "tag-for-s3", "Environment":"prod"},
  {"Owner": "bridgecrew"
  },
  {"yor_trace": "4329587194",
    "git_org": "bana"})
}