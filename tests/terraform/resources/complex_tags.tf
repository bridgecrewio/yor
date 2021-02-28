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

resource "aws_autoscaling_group" "aurora_cluster_bastion_auto_scaling_group" {
  default_cooldown          = "300"
  desired_capacity          = "1"
  force_delete              = "false"
  health_check_grace_period = "60"
  health_check_type         = "EC2"
  max_instance_lifetime     = "0"
  max_size                  = "1"
  metrics_granularity       = "1Minute"
  min_size                  = "1"
  name                      = "bc-aurora-cluster-bastion-auto-scaling-group"
  protect_from_scale_in     = "false"

  wait_for_capacity_timeout = "10m"


  tag {
    key                 = "Name"
    propagate_at_launch = true
    value               = "bc-aurora-bastion"
  }

  tag {
    key = "Env"
    propagate_at_launch = false
    value = "prod"
  }

  tags = {
    git_org              = "bridgecrewio"
    git_repo             = "platform"
    yor_trace            = "48564943-4cfc-403c-88cd-cbb207e0d33e"
  }
}