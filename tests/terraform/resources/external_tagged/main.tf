resource "aws_s3_bucket" "a" {
  bucket = "my-tf-test-bucket"
  acl = "private"

  tags = {
    Name = "My bucket"
  }
}

resource "aws_s3_bucket" "b" {
  bucket = "my-tf-test-bucket"
  acl = "private"

  tags = {
    Name = "My other bucket"
  }
}


resource "aws_s3_bucket" "c" {
  bucket = "my-tf-test-bucket"
  acl = "private"
}

resource "aws_vpc" "d" {
  cidr_block = ""
}

