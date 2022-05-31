resource "aws_s3_bucket" "bucket_module3" {
  bucket = "tf-test-bucket3"
  acl    = "public_read"
  versioning {
    enabled = false
  }

  tags = {
    "Name"    = "bucket3"
    yor_trace = "887ca272-5af1-4cbb-b7d6-9f3ba50842a6"
  }
}
