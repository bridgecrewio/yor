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
  tags = {}
}
