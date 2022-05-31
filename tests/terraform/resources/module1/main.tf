resource "aws_s3_bucket" "bucket_module1" {
  bucket = "tf-test-bucket-destination-12345"
  acl    = "private"
  versioning {
    enabled = false
  }
  tags = {
    yor_trace = "67da8fbd-99ba-42e7-8924-91c37ee02589"
  }
}

module "moduleRef" {
  source = "../module2"
}
