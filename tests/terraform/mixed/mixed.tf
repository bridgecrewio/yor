resource aws_s3_bucket "test-bucket" {
  name = "test-bucket"
  tags = {
    git_commit           = "N/A"
    git_file             = "tests/terraform/mixed/mixed.tf"
    git_last_modified_at = "2021-04-01 11:30:10"
    git_last_modified_by = "james.woolfenden@gmail.com"
    git_modifiers        = "james.woolfenden/nimrodkor"
    git_org              = "JamesWoolfenden"
    git_repo             = "yor"
    yor_trace            = "0c25f5bb-8298-4ff3-9824-87b941b44e81"
  }
}

resource tls_private_key "pem" {
  algorithm = "some-algo"
}