provider "aws" {
  region = "us-east-1"
  profile = "dev2"
  version = "3.27"
}

resource "aws_s3_bucket" "bucket_with_var" {
  bucket = "tf-test-bucket-destination-12345"
  acl = "private"
  versioning {
    enabled = false
  }

  tags = var.bucket_tags
}
