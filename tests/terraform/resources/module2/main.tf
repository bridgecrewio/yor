resource "aws_s3_bucket" "bucket_module2" {
  bucket = "tf-test-bucket2"
  acl = "public_read"
  versioning {
    enabled = false
  }

  tags = {
    "Name" = "bucket2"
  }
}
