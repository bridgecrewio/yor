resource "aws_instance" "many_instance_tags" {
  ami           = ""
  instance_type = ""
  tags = merge({ "Name" = "tag-for-instance", "Environment" = "prod" },
    { "Owner" = "bridgecrew"
    },
    { "yor_trace"          = "4329587194",
      "git_org"            = "" }, {
      git_commit           = "0000000000000000000000000000000000000000"
      git_file             = "....teststerraformresourcescomplex_tags.tf"
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
      git_file             = "....teststerraformresourcescomplex_tags.tf"
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
    git_file             = "....teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "4f9373f9-3ba8-461a-b067-189b20bfc8eb"
  })
}

resource "aws_s3_bucket" "bucket_var_tags" {
  tags = merge(var.tags, {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "....teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "bae7f24f-45f4-436a-b98a-aaeed6e4aff3"
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
      git_file             = "....teststerraformresourcescomplex_tags.tf"
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
    git_file             = "....teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "1861aa93-ae51-40f6-9e5e-ccbadc88cde1"
  })
}

resource "aws_instance" "instance_empty_tag" {
  ami           = ""
  instance_type = ""
  tags = merge({}, {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "....teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "f57b258b-b562-4166-9976-467c5879705a"
  })
}

resource "aws_instance" "instance_no_tags" {
  ami           = ""
  instance_type = ""
  tags = {
    git_commit           = "0000000000000000000000000000000000000000"
    git_file             = "....teststerraformresourcescomplex_tags.tf"
    git_last_modified_at = "2020-06-16 17:46:24"
    git_last_modified_by = "user@gmail.com"
    git_modifiers        = "user"
    git_org              = ""
    git_repo             = ""
    yor_trace            = "b55edc53-6322-445f-a983-45d5000ee620"
  }
}