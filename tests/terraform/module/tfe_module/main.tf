module "tfe_module" {
  source = "app.terraform.io/path/to/module/aws"
  tags = {
    Application = "application"
    Env         = var.env
    yor_trace   = "fd59e45a-5785-41b1-8195-54cf9971fd47"
  }
}