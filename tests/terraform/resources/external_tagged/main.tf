resource "aws_s3_bucket" "a" {
  bucket = "my-tf-test-bucket"
  acl = "private"

  tags = {
    Name = "My bucket"
    env = "dev"
    git_commit = "00193660c248483862c06e2ae96111adfcb683af"
    git_file = "tests/terraform/resources/external_tagged/main.tf"
    git_last_modified_at = "2021-06-01 07:58:29"
    git_last_modified_by = "tron47@gmail.com"
    git_modifiers = "tron47"
    git_org = "bridgecrewio"
    git_repo = "yor"
    team = "seceng"
    yor_trace = "fcf8717c-372b-4d9b-9dab-78f5af48fc56"
  }
}

resource "aws_s3_bucket" "c" {
  bucket = "my-tf-test-bucket"
  acl = "private"
}

resource "aws_vpc" "d" {
  cidr_block = ""
  tags = {
    git_commit = "00193660c248483862c06e2ae96111adfcb683af"
    git_file = "tests/terraform/resources/external_tagged/main.tf"
    git_last_modified_at = "2021-06-01 07:58:29"
    git_last_modified_by = "tron47@gmail.com"
    git_modifiers = "tron47"
    git_org = "bridgecrewio"
    git_repo = "yor"
    team = "seceng"
    yor_trace = "181119e5-c2ff-4ba1-a35e-78d2c0793297"
  }
}

