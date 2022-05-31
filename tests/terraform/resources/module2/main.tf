resource "aws_s3_bucket" "bucket_module2" {
  bucket = "tf-test-bucket2"
  acl    = "public_read"
  versioning {
    enabled = false
  }

  tags = {
    "Name"    = "bucket2"
    yor_trace = "1d66c4c6-f2bd-4a25-a64c-eabbef0e5285"
  }
}

module "moduleRef2" {
  source = "../module3"
}
