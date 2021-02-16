resource "aws_s3_bucket" "bucket_module3" {
  bucket = "tf-test-bucket3"
  acl = "public_read"
  versioning {
    enabled = false
  }

  tags = {
    "Name" = "bucket3"
  }
}
