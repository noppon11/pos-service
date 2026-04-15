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

    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
  }
}

provider "null" {}

provider "github" {
  owner = "noppon11"
  token = var.github_token
}

resource "null_resource" "example" {
  triggers = {
    app_name    = var.app_name
    environment = var.environment
  }
}

# ลบ github_repository_ruleset ออก แล้วใช้อันนี้แทน
resource "github_branch_protection" "main" {
  repository_id = "pos-service"
  pattern       = "main"

  required_status_checks {
    strict   = true
    contexts = ["go-test"]
  }

  required_pull_request_reviews {
    required_approving_review_count = 1
    dismiss_stale_reviews           = true
    require_last_push_approval      = false
  }

  enforce_admins = false
}