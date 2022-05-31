resource aws_s3_bucket "test-bucket" {
  name = "test-bucket"
  tags = {
    yor_trace = "19156f28-6814-4ae5-a3b1-440132b749c7"
  }
}

resource tls_private_key "pem" {
  algorithm = "some-algo"
}