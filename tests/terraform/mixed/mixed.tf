resource aws_s3_bucket "test-bucket" {
  name = "test-bucket"
}

resource tls_private_key "pem" {
  algorithm = "some-algo"
}