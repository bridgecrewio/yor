resource "google_monitoring_notification_channel" "this" {
  display_name = "GDOS Development Team email"
  force_delete = false
  labels = {
    email_address = "fake_email@blahblah.com"
  }
  project = module.project.project_id
  type    = "email"
}