module "tfe_module" {
  source  = "app.terraform.io/path/to/module/aws"
  tags = {
    Application = "application"
    Env         = var.env
  }
}