variable "scw_project_id" {
  type        = string
  description = "Scaleway Project ID (e.g. from `scw config get default-project-id`)."
}

variable "scw_zone" {
  type        = string
  default     = "fr-par-1"
  description = "Scaleway zone. Serverless Containers and Serverless SQL DB are available in the Paris region."
}

variable "db_name" {
  type        = string
  default     = "tulospalvelu"
  description = "Name of the Serverless SQL Database."
}

variable "db_max_cpu" {
  type        = number
  default     = 8
  description = "Maximum vCPUs for the Serverless SQL Database autoscaling."
}

variable "deploy_container" {
  type        = bool
  default     = false
  description = "Deploy the container on apply. Keep false for the bootstrap run (the registry image must be pushed first), then set true."
}

variable "container_image_tag" {
  type        = string
  default     = "latest"
  description = "Tag of the image pushed to the Scaleway Container Registry."
}
