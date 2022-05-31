provider "aws" {

}

resource "aws_s3_bucket" "financials" {
  # bucket is not encrypted
  # bucket does not have access logs
  # bucket does not have versioning
  bucket        = "yor-test-1"
  acl           = "private"
  force_destroy = true

  tags = {
    yor_trace = "b123cc2d-026e-4de4-82c4-66a23692e859"
  }
}