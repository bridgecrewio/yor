resource "aws_kms_key" "logs_key" {
  # key does not have rotation enabled
  description = "modified_description"

  deletion_window_in_days = 7
  tags = {
    yor_trace = "c5bca344-27d6-484d-9b57-87cfd66da7fa"
  }
}

resource "aws_kms_alias" "logs_key_alias" {
  name          = "alias/${local.resource_prefix.value}-logs-bucket-key"
  target_key_id = "${aws_kms_key.logs_key.key_id}"
}
