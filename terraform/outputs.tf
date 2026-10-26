output "container_url" {
  value     = scaleway_container.app.public_endpoint
  sensitive = true
}

output "registry_endpoint" {
  value = scaleway_registry_namespace.image.endpoint
}

output "registry_image" {
  value = "${scaleway_registry_namespace.image.endpoint}/tulospalvelu:${var.container_image_tag}"
}

output "db_endpoint" {
  value     = scaleway_sdb_sql_database.db.endpoint
  sensitive = true
}

output "db_connection_string" {
  description = "psql-ready connection string (IAM user / API key secret)."
  value       = local.database_url
  sensitive   = true
}
