# Tulospalvelu: Serverless Container + Serverless SQL Database.
#
# The database is reached over its public endpoint with IAM
# authentication: DB username = IAM application ID, DB password = IAM
# API key secret. TLS is required (sslmode=require).

resource "scaleway_registry_namespace" "image" {
  name        = "tulospalvelu"
  description = "Container Registry namespace for tulospalvelu images"
}

resource "scaleway_container_namespace" "app" {
  name        = "tulospalvelu"
  description = "Serverless Containers namespace for tulospalvelu"
}

resource "scaleway_sdb_sql_database" "db" {
  name    = var.db_name
  min_cpu = 0
  max_cpu = var.db_max_cpu
}

# IAM identity used by the container to authenticate to the database.
resource "scaleway_iam_application" "app" {
  name = "tulospalvelu-app"
}

resource "scaleway_iam_api_key" "app" {
  application_id = scaleway_iam_application.app.id
}

resource "scaleway_iam_policy" "db_access" {
  name           = "tulospalvelu-db-access"
  description    = "Grants the tulospalvelu container read/write access to the Serverless SQL Database"
  application_id = scaleway_iam_application.app.id

  rule {
    project_ids          = [var.scw_project_id]
    permission_set_names = ["ServerlessSQLDatabaseReadWrite"]
  }
}

locals {
  # sdb endpoint looks like:
  #   postgres://<uuid>.pg.sdb.fr-par.scw.cloud:5432/<name>?sslmode=require
  db_host = trimsuffix(trimprefix(regex("://[^:]+", scaleway_sdb_sql_database.db.endpoint), "://"), ":")
  db_port = trimprefix(regex(":[0-9]{1,5}", scaleway_sdb_sql_database.db.endpoint), ":")
  db_name = trimsuffix(element(split("/", scaleway_sdb_sql_database.db.endpoint), 3), "?sslmode=require")

  database_url = format("postgres://%s:%s@%s:%s/%s?sslmode=require",
    scaleway_iam_application.app.id,
    scaleway_iam_api_key.app.secret_key,
    local.db_host,
    local.db_port,
    local.db_name,
  )

  registry_image = "${scaleway_registry_namespace.image.endpoint}/tulospalvelu:${var.container_image_tag}"
}

resource "scaleway_container" "app" {
  name         = "tulospalvelu"
  description  = "Hiihtokisojen tulospalvelu (Go)"
  namespace_id = scaleway_container_namespace.app.id

  registry_image = local.registry_image
  port           = 8080

  cpu_limit          = 1000
  memory_limit_bytes = 536870912 # 512 MiB

  min_scale = 0
  max_scale = 2
  timeout   = 20

  protocol = "http1"
  privacy  = "public"
  deploy   = var.deploy_container

  environment_variables = {
    GEO_RESTRICT = "true"
    TRUST_XFF    = "true"
  }

  secret_environment_variables = {
    DATABASE_URL = local.database_url
  }
}
