resource "aws_instance" "many_instance_tags" {
  ami           = ""
  instance_type = ""
  tags = merge({ "Name" = "tag-for-instance", "Environment" = "prod" },
    { "Owner" = "bridgecrew"
    },
    { "yor_trace"          = "4329587194",
      "git_org"            = "bridgecrewio" }, {
      git_commit           = "0000000000000000000000000000000000000000"
      git_file             = "tests/terraform/resources/complex_tags.tf"
      git_last_modified_at = "2020-06-16 17:46:24"
      git_last_modified_by = "user@gmail.com"
      git_modifiers        = "user"
      git_repo             = "yor"
  })
}

resource "aws_alb" "alb_with_merged_tags" {
  tags = merge({ "Name" = "tag-for-alb", "Environment" = "prod" },
    { "yor_trace"          = "4329587194",
      "git_org"            = "bridgecrewio" }, {
      git_commit           = "0000000000000000000000000000000000000000"
      git_file             = "tests/terraform/resources/complex_tags.tf"
      git_last_modified_at = "2020-06-16 17:46:24"
      git_last_modified_by = "user@gmail.com"
      git_modifiers        = "user"
      git_repo             = "yor"
  })
}

resource "aws_vpc" "vpc_tags_one_line" {
  cidr_block = ""
  tags = merge({ "Name" = "tag-for-s3", "Environment" = "prod" }, {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "tests/terraform/resources/complex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = "bridgecrewio"
    git_repo             = "yor"
    yor_trace            = "200a28f3-5977-439e-b802-1ac3ca033d4b"
  })
}

resource "aws_s3_bucket" "bucket_var_tags" {
  tags = merge(var.tags, {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "tests/terraform/resources/complex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = "bridgecrewio"
    git_repo             = "yor"
    yor_trace            = "b2b3a7a6-6aed-4148-a7d0-cc0003ca0256"
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
      git_commit           = "0000000000000000000000000000000000000000"
      git_file             = "tests/terraform/resources/complex_tags.tf"
      git_last_modified_at = "2020-06-16 17:46:24"
      git_last_modified_by = "user@gmail.com"
      git_modifiers        = "user"
      git_repo             = "yor"
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
    git_file             = "tests/terraform/resources/complex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = "bridgecrewio"
    git_repo             = "yor"
    yor_trace            = "ad24ee3a-a7e7-44d2-b89b-52b41240c6c0"
  })
}

resource "aws_instance" "instance_empty_tag" {
  ami           = ""
  instance_type = ""
  tags = merge({}, {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "tests/terraform/resources/complex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = "bridgecrewio"
    git_repo             = "yor"
    yor_trace            = "ea52b61f-bb51-4fbc-9ad1-15cd710da64f"
  })
}

resource "aws_instance" "instance_no_tags" {
  ami           = ""
  instance_type = ""
  tags = {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "tests/terraform/resources/complex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = "bridgecrewio"
    git_repo             = "yor"
    yor_trace            = "0584a9b1-c057-485f-aab2-8017130a0ea6"
  }
}