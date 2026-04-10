terraform {
  required_version = ">= 1.5.0"

  cloud {
    organization = "pp-aura-wellness"

    workspaces {
      name = "pos-service-prod"
    }
  }

  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.2"
    }
  }
}

provider "null" {}

resource "null_resource" "example" {
  triggers = {
    app_name    = var.app_name
    environment = var.environment
  }
}