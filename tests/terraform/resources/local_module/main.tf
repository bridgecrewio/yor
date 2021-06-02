module "sub_module" {
  source =   "./sub_local_module"
  tags = {
    Name = "test"
  }
}