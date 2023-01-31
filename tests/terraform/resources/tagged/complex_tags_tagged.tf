resource "aws_instance" "many_instance_tags" {
  ami           = ""
  instance_type = ""
  tags = merge({ "Name" = "tag-for-instance", "Environment" = "prod" },
    { "Owner" = "bridgecrew"
    },
    { yor_trace            = "4329587194",
      git_org              = "bridgecrewio" }, {
      git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
      git_file             = "README.md"
      git_last_modified_at = "2020-03-28 21:42:46"
      git_last_modified_by = "schosterbarak@gmail.com"
      git_modifiers        = "jonjozwiak/schosterbarak"
      git_repo             = "terragoat"
  })
}

resource "aws_alb" "alb_with_merged_tags" {
  tags = merge({ "Name" = "tag-for-alb", "Environment" = "prod" },
    { yor_trace            = "4329587194",
      git_org              = "bridgecrewio" }, {
      git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
      git_file             = "README.md"
      git_last_modified_at = "2020-03-28 21:42:46"
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
    git_last_modified_at = "2020-03-28 21:42:46"
    git_last_modified_by = "schosterbarak@gmail.com"
    git_modifiers        = "jonjozwiak/schosterbarak"
    git_org              = "bridgecrewio"
    git_repo             = "terragoat"
    yor_trace            = "85da7acc-b505-49e4-8f42-bb5e708c1aa3"
  })
}

resource "aws_s3_bucket" "bucket_var_tags" {
  tags = merge(var.tags, {
    git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
    git_file             = "README.md"
    git_last_modified_at = "2020-03-28 21:42:46"
    git_last_modified_by = "schosterbarak@gmail.com"
    git_modifiers        = "jonjozwiak/schosterbarak"
    git_org              = "bridgecrewio"
    git_repo             = "terragoat"
    yor_trace            = "a7698353-b81c-4ddb-bb9f-745718f8c7ae"
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
    { yor_trace            = "4329587194",
      git_org              = "bridgecrewio" }, {
      git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
      git_file             = "README.md"
      git_last_modified_at = "2020-03-28 21:42:46"
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
    git_last_modified_at = "2020-03-28 21:42:46"
    git_last_modified_by = "schosterbarak@gmail.com"
    git_modifiers        = "jonjozwiak/schosterbarak"
    git_org              = "bridgecrewio"
    git_repo             = "terragoat"
    yor_trace            = "a1cb42d1-bfbb-486e-8d79-31bedc19c293"
  })
}

resource "aws_instance" "instance_empty_tag" {
  ami           = ""
  instance_type = ""
  tags = merge({}, {
    git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
    git_file             = "README.md"
    git_last_modified_at = "2020-03-28 21:42:46"
    git_last_modified_by = "schosterbarak@gmail.com"
    git_modifiers        = "jonjozwiak/schosterbarak"
    git_org              = "bridgecrewio"
    git_repo             = "terragoat"
    yor_trace            = "8ec7f549-4133-4dcc-bdb9-86ab0f336d9c"
  })
}

resource "aws_instance" "instance_no_tags" {
  ami           = ""
  instance_type = ""
  tags = {
    git_commit           = "47accf06f13b503f3bab06fed7860e72f7523cac"
    git_file             = "README.md"
    git_last_modified_at = "2020-03-28 21:42:46"
    git_last_modified_by = "schosterbarak@gmail.com"
    git_modifiers        = "jonjozwiak/schosterbarak"
    git_org              = "bridgecrewio"
    git_repo             = "terragoat"
    yor_trace            = "a51f6e65-cd2d-4f53-962c-0d2894fc6418"
  }
}
