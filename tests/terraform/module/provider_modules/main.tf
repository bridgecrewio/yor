module "project-factory" {
  source  = "terraform-google-modules/project-factory/google"
  version = "11.0.0"
  labels = {
    test      = "true"
    yor_trace = "06d4ac98-fa1b-4fa7-af17-ec7938c4ba53"
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "3.2.0"
  tags = {
    test      = "true"
    yor_trace = "e4c69f45-dd72-43b6-a9bb-63094827d69a"
  }
}

module "project-factory_github" {
  source = "github.com/terraform-google-modules/terraform-google-project-factory"
  labels = {
    test      = "true"
    yor_trace = "6fdff67e-8970-4317-b17c-eceeefc7ea6b"
  }
}

module "project-factory_git" {
  source = "git@github.com:terraform-google-modules/terraform-google-project-factory.git"
  labels = {
    test      = "true"
    yor_trace = "dd81040a-1d16-4216-afcf-4be2e9ee7a57"
  }
}

module "caf" {
  source = "aztfmod/caf/azurerm"
  tags = {
    test      = "true"
    yor_trace = "10738839-d087-4f60-b5e5-555ac268a05d"
  }
}

module "caf" {
  source = "git@github.com:aztfmod/terraform-azurerm-caf.git"
  tags = {
    test      = "true"
    yor_trace = "10738839-d087-4f60-b5e5-555ac268a05d"
  }
}

module "bastion" {
  source = "oracle-terraform-modules/bastion/oci"
  freeform_tags = {
    test      = "true"
    yor_trace = "1f98042b-6c2d-4f12-b225-5a65b3f54426"
  }
}

module "run-common_logs" {
  // Tags attribute is extra_tags
  source  = "claranet/run-common/azurerm//modules/logs"
  version = "3.0.0"
  extra_tags = {
    test      = "true"
    yor_trace = "f318b137-9300-4051-9585-928113c76c13"
  }
}