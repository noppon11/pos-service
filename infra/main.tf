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

resource "github_repository_ruleset" "main_require_tests" {
  name        = "require-tests-before-merge"
  repository  = "pos-service"
  target      = "branch"
  enforcement = "active"

  conditions {
    ref_name {
      include = ["~DEFAULT_BRANCH"]
      exclude = []
    }
  }

  rules {
    pull_request {
      dismiss_stale_reviews_on_push     = false
      require_code_owner_review         = false
      require_last_push_approval        = false
      required_approving_review_count   = 1
      required_review_thread_resolution = true
    }

    required_status_checks {
      strict_required_status_checks_policy = true

      required_check {
        context = "go-test"
      }
    }
  }
}