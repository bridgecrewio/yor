resource "aws_s3_bucket" "bucket_module1" {
  bucket = "tf-test-bucket-destination-12345"
  acl = "private"
  versioning {
    enabled = false
  }
}

module "moduleRef" {
  source = "../module2"
}