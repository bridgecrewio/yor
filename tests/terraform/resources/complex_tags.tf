resource "aws_instance" "many_instance_tags" {
  ami           = ""
  instance_type = ""
  tags = merge({ "Name" = "tag-for-instance", "Environment" = "prod" },
    { "Owner" = "bridgecrew"
    },
    { "yor_trace" = "4329587194",
      "git_org"   = "bana" }, {
      Name        = "tag-for-instance"
  })
}

resource "aws_alb" "alb_with_merged_tags" {
  tags = merge({ "Name" = "tag-for-alb", "Environment" = "prod" },
    { "yor_trace" = "4329587194",
      "git_org"   = "bana" }, {
      Name        = "tag-for-alb"
  })
}

resource "aws_vpc" "vpc_tags_one_line" {
  cidr_block = ""
  tags = { "Name" = "tag-for-s3", "Environment" = "prod"
    yor_trace = "aa5aedf1-f46f-409c-91ba-1c712b7880d9"
  }
}

resource "aws_s3_bucket" "bucket_var_tags" {
  tags = merge(var.tags, {
    yor_trace = "28941e6c-06b4-4901-923b-7477080b2519"
  })
}

variable "tags" {
  default = {
    "Name"        = "tag-for-s3"
    "Environment" = "prod"
  }
}

resource "aws_instance" "instance_merged_var" {
  ami           = ""
  instance_type = ""
  tags = merge(var.tags,
    { "yor_trace" = "4329587194",
  "git_org" = "bana" })
}

variable "new_env_tag" {
  default = {
    "Environment" = "old_env"
  }
}

resource "aws_instance" "instance_merged_override" {
  ami           = ""
  instance_type = ""
  tags = merge(var.new_env_tag, { "Environment" = "new_env" }, {
    yor_trace = "3b8e0da3-a795-44ec-95c7-d80256bcdea5"
  })
}

resource "aws_instance" "instance_empty_tag" {
  ami           = ""
  instance_type = ""
  tags = {
    yor_trace = "7938e439-cb58-4fed-9024-d957e4f5420b"
  }
}

resource "aws_instance" "instance_no_tags" {
  ami           = ""
  instance_type = ""
  tags = {
    yor_trace = "86b90701-e0dd-487b-b009-df18464e033a"
  }
}

resource "aws_instance" "instance_null_tags" {
  ami           = ""
  instance_type = ""
  tags = {
    yor_trace = "16dccc6d-81d7-4a06-ae9b-2c37580031c3"
  }
}

resource "aws_autoscaling_group" "autoscaling_group_tagged" {
  // This resource should not be tagged
  tag {
    key                 = "Name"
    propagate_at_launch = false
    value               = "Mine"
  }
  max_size = 0
  min_size = 0
}

resource "aws_autoscaling_group" "autoscaling_group" {
  // This resource should not be tagged as well
  max_size = 0
  min_size = 0
}