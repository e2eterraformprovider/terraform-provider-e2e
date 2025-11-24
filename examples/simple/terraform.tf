terraform {
  required_providers {
    e2e = {
      source  = "registry.terraform.io/e2eterraformprovider/e2e"
      version = "0.1.0"
    }
  }
}

provider "e2e" {
  api_key    = var.api_key
  auth_token = var.auth_token
  # api_endpoint is optional, defaults to https://api.e2enetworks.com/myaccount/api/v1/
}

