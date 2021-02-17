resource "aws_instance" "many_instance_tags" {
  ami           = ""
  instance_type = ""
  tags = merge({ "Name" = "tag-for-instance", "Environment" = "prod" },
    { "Owner" = "bridgecrew"
    },
    { "yor_trace"          = "4329587194",
      "git_org"            = "bridgecrewio" }, {
      git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
      git_file             = "README.md"
      git_last_modified_at = "2020-03-28 21:42:46 +0000 UTC"
      git_last_modified_by = "schosterbarak@gmail.com"
      git_modifiers        = "jonjozwiak/schosterbarak"
      git_repo             = "terragoat"
  })
}

resource "aws_alb" "alb_with_merged_tags" {
  tags = merge({ "Name" = "tag-for-alb", "Environment" = "prod" },
    { "yor_trace"          = "4329587194",
      "git_org"            = "bridgecrewio" }, {
      git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
      git_file             = "README.md"
      git_last_modified_at = "2020-03-28 21:42:46 +0000 UTC"
      git_last_modified_by = "schosterbarak@gmail.com"
      git_modifiers        = "jonjozwiak/schosterbarak"
      git_repo             = "terragoat"
  })
}

resource "aws_vpc" "vpc_tags_one_line" {
  cidr_block = ""
  tags = merge({ "Name" = "tag-for-s3", "Environment" = "prod" }, {
    git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
    git_file             = "README.md"
    git_last_modified_at = "2020-03-28 21:42:46 +0000 UTC"
    git_last_modified_by = "schosterbarak@gmail.com"
    git_modifiers        = "jonjozwiak/schosterbarak"
    git_org              = "bridgecrewio"
    git_repo             = "terragoat"
    yor_trace            = "17dc538c-bbed-488e-b5cf-b8eb1e433dc7"
  })
}

resource "aws_s3_bucket" "bucket_var_tags" {
  tags = merge(var.tags, {
    git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
    git_file             = "README.md"
    git_last_modified_at = "2020-03-28 21:42:46 +0000 UTC"
    git_last_modified_by = "schosterbarak@gmail.com"
    git_modifiers        = "jonjozwiak/schosterbarak"
    git_org              = "bridgecrewio"
    git_repo             = "terragoat"
    yor_trace            = "17dc538c-bbed-488e-b5cf-b8eb1e433dc7"
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
    { "yor_trace"          = "4329587194",
      "git_org"            = "bridgecrewio" }, {
      git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
      git_file             = "README.md"
      git_last_modified_at = "2020-03-28 21:42:46 +0000 UTC"
      git_last_modified_by = "schosterbarak@gmail.com"
      git_modifiers        = "jonjozwiak/schosterbarak"
      git_repo             = "terragoat"
  })
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
    git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
    git_file             = "README.md"
    git_last_modified_at = "2020-03-28 21:42:46 +0000 UTC"
    git_last_modified_by = "schosterbarak@gmail.com"
    git_modifiers        = "jonjozwiak/schosterbarak"
    git_org              = "bridgecrewio"
    git_repo             = "terragoat"
    yor_trace            = "17dc538c-bbed-488e-b5cf-b8eb1e433dc7"
  })
}