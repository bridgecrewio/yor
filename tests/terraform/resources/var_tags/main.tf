provider "aws" {
  region  = "us-east-1"
  profile = "dev2"
  version = "3.27"
}

resource "aws_s3_bucket" "bucket_with_var" {
  bucket = "tf-test-bucket-destination-12345"
  acl    = "private"
  versioning {
    enabled = false
  }

  tags = merge(var.bucket_tags, {
    yor_trace = "e8ae4ed9-4150-4158-8520-bc735c0ab8e0"
  })
}
