resource "aws_instance" "many_instance_tags" {
  ami = ""
  instance_type = ""
  tags =  merge({"Name" = "tag-for-instance", "Environment" = "prod"},
  {"Owner" = "bridgecrew"
  },
  {"yor_trace" = "4329587194",
    "git_org" = "bana"})
}

resource "aws_alb" "alb_with_merged_tags" {
  tags = merge({"Name" = "tag-for-alb", "Environment" = "prod"},
  {"yor_trace" = "4329587194",
    "git_org" = "bana"})
}

resource "aws_vpc" "vpc_tags_one_line" {
  cidr_block = ""
  tags = {"Name" = "tag-for-s3", "Environment" = "prod"}
}

resource "aws_s3_bucket" "bucket_var_tags" {
  tags = var.tags
}

variable "tags" {
  default = {
    "Name" = "tag-for-s3"
    "Environment" ="prod"
  }
}

resource "aws_instance" "instance_merged_var" {
  ami = ""
  instance_type = ""
  tags =  merge(var.tags,
  {"yor_trace" = "4329587194",
    "git_org" = "bana"})
}

variable "new_env_tag" {
  default = {
    "Environment" = "old_env"
  }
}

resource "aws_instance" "instance_merged_override" {
  ami = ""
  instance_type = ""
  tags = merge(var.new_env_tag, {"Environment" = "new_env"})
}

resource "aws_instance" "instance_empty_tag" {
  ami = ""
  instance_type = ""
  tags = {}
}

resource "aws_instance" "instance_no_tags" {
  ami = ""
  instance_type = ""
}

resource "aws_instance" "instance_null_tags" {
  ami = ""
  instance_type = ""
  tags = null
}

resource "aws_autoscaling_group" "autoscaling_group_tagged" {
  // This resource should not be tagged
  tag {
    key = "Name"
    propagate_at_launch = false
    value = "Mine"
  }
  max_size = 0
  min_size = 0
}

resource "aws_autoscaling_group" "autoscaling_group" {
  // This resource should not be tagged as well
  max_size = 0
  min_size = 0
}