module "network" {
  source = "Azure/network/azurerm"
  tags = {
    test = "true"
    yor_toggle = false
  }
}
