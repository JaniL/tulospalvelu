terraform {
  required_version = ">= 1.5"

  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = "~> 2.0"
    }
  }
}

provider "scaleway" {
  zone = var.scw_zone
  # Credentials come from standard SCW_ACCESS_KEY / SCW_SECRET_KEY
  # environment variables (or `scw` CLI profile / SCW_CONFIG_FILE).
  # project_id is set explicitly so all resources land in the right project.
  project_id = var.scw_project_id
}
