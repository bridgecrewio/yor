resource "aws_instance" "many_instance_tags" {
  ami           = ""
  instance_type = ""
  tags = merge({ "Name" = "tag-for-instance", "Environment" = "prod" },
    { "Owner" = "bridgecrew"
    },
    { "yor_trace"          = "4329587194",
      "git_org"            = "" }, {
      git_commit           = "0000000000000000000000000000000000000000"
      git_file             = "......teststerraformresourcescomplex_tags.tf"
      git_last_modified_at = "2020-06-16 17:46:24"
      git_last_modified_by = "user@gmail.com"
      git_modifiers        = "user"
      git_repo             = ""
  })
}

resource "aws_alb" "alb_with_merged_tags" {
  tags = merge({ "Name" = "tag-for-alb", "Environment" = "prod" },
    { "yor_trace"          = "4329587194",
      "git_org"            = "" }, {
      git_commit           = "0000000000000000000000000000000000000000"
      git_file             = "......teststerraformresourcescomplex_tags.tf"
      git_last_modified_at = "2020-06-16 17:46:24"
      git_last_modified_by = "user@gmail.com"
      git_modifiers        = "user"
      git_repo             = ""
  })
}

resource "aws_vpc" "vpc_tags_one_line" {
  cidr_block = ""
  tags = merge({ "Name" = "tag-for-s3", "Environment" = "prod" }, {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "......teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "47127974-d01d-4ec0-b2b1-231e96096be5"
  })
}

resource "aws_s3_bucket" "bucket_var_tags" {
  tags = merge(var.tags, {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "......teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "8f9bd7da-aadf-49f8-a12f-59347941d699"
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
      "git_org"            = "" }, {
      git_commit           = "0000000000000000000000000000000000000000"
      git_file             = "......teststerraformresourcescomplex_tags.tf"
      git_last_modified_at = "2020-06-16 17:46:24"
      git_last_modified_by = "user@gmail.com"
      git_modifiers        = "user"
      git_repo             = ""
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
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "......teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "fddc0ef6-ae83-4344-9f80-84cb5ada0693"
  })
}

resource "aws_instance" "instance_empty_tag" {
  ami           = ""
  instance_type = ""
  tags = merge({}, {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "......teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "d8ec2406-0064-45b3-8a70-b06c66597d9f"
  })
}

resource "aws_instance" "instance_no_tags" {
  ami           = ""
  instance_type = ""
  tags = {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "......teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "6fca1c11-3a33-4dbf-b8a7-5f8e59bac1e2"
  }
}