resource "google_storage_bucket" "auto-expire" {
  name          = "yor-test-gcs-bucket"
  location      = "US"
  force_destroy = true

  lifecycle_rule {
    condition {
      age = 3
    }
    action {
      type = "Delete"
    }
  }
}