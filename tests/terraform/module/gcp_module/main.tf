module "project" {
  source = "git@github.com:my-org/terraform-google-project?ref=v2.1.3"

  cost_center       = "uZZZ"
  folder_id         = var.folder_id
  prefix            = "shared"
  system            = "services"
  env               = var.env
  random_project_id = var.random_project_id

  labels = {
    team                 = "ops"
    yor_trace            = "8019ee39-ca38-4a8d-8b7a-d50ffd353f5c"
    git_commit           = "c828864803f716011793b240628c4311330de7d8"
    git_file             = "global/main.tf"
    git_last_modified_at = "2021-11-03 13:13:21"
    git_last_modified_by = "my-name@my-org.com"
    git_modifiers        = "my-name"
    git_org              = "my-org"
    git_repo             = "my-gcp-repo"
  }
}