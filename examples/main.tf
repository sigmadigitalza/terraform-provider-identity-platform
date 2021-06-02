terraform {
  required_providers {
    identity-platform = {
      version = "1.0.0"
      source = "sigmadigital.io/google/identity-platform"
    }
  }
}

provider "identity-platform" {}

resource "identity_platform_config" "auth_config" {
  provider = identity-platform

  project_id = "<gcp-project-id>"

  email {
    enabled = true
    password_required = true
  }

  phone_number {
    enabled = true
  }

  notification {
    send_email {
      method = "default"
      callback_uri = "https://<gcp-project-id>.firebaseapp.com/__/auth/action"
    }
  }

  subtype = "IDENTITY_PLATFORM"

  authorized_domains = [
    "localhost",
    "<gcp-project-id>.firebaseapp.com",
    "<gcp-project-id>.web.app",
  ]
}
