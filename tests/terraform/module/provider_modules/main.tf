module "project-factory" {
  source  = "terraform-google-modules/project-factory/google"
  version = "11.0.0"
  labels  = {
    test = "true"
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "3.2.0"
  tags    = {
    test = "true"
  }
}

module "project-factory_github" {
  source = "github.com/terraform-google-modules/terraform-google-project-factory"
  labels = {
    test = "true"
  }
}

module "project-factory_git" {
  source = "git@github.com:terraform-google-modules/terraform-google-project-factory.git"
  labels = {
    test = "true"
  }
}

module "caf" {
  source = "aztfmod/caf/azurerm"
  tags   = {
    test = "true"
  }
}

module "caf" {
  source = "git@github.com:aztfmod/terraform-azurerm-caf.git"
  tags   = {
    test = "true"
  }
}

module "bastion" {
  source        = "oracle-terraform-modules/bastion/oci"
  freeform_tags = {
    test = "true"
  }
}

module "run-common_logs" {
  // Tags attribute is extra_tags
  source  = "claranet/run-common/azurerm//modules/logs"
  version = "3.0.0"
  extra_tags = {
    test = "true"
  }
}