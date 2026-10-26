# Deploying tulospalvelu (Go) to Scaleway

The deployment target is:

- **Serverless Container** running the Go image (scales to zero)
- **Serverless SQL Database** (PostgreSQL, scales to zero, IAM auth)
- **Container Registry** namespace for the image

Terraform creates everything. The image is built and pushed with
**podman** (no Docker required). Everything below runs on the machine
that has Scaleway access.

## Prerequisites

1. **Scaleway credentials** — an API key with (at minimum) these
   permission sets for the operator:
   `ContainersFullAccess`, `ContainersRegistryFullAccess`,
   `ServerlessSQLDatabaseFullAccess`, plus IAM permissions to create
   applications/API keys/policies.

   ```bash
   export SCW_ACCESS_KEY=...
   export SCW_SECRET_KEY=...
   ```

   (Or log in with `scw` CLI and Terraform picks up its profile.)

2. **Project ID**:

   ```bash
   export SCW_PROJECT_ID=$(scw config get default-project-id)
   ```

3. `terraform` and `podman` installed.

## One-time setup

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# edit scw_project_id in terraform.tfvars
```

## Deploy

### 1. Bootstrap (creates DB, registry, namespaces — no container yet)

The container cannot be deployed before the image exists in the
registry, so the first apply runs with `deploy_container = false`
(default):

```bash
terraform init
terraform apply
```

### 2. Build and push the image

```bash
cd ..
make image REGISTRY=$(cd terraform && terraform output -raw registry_image)
# or manually:
# podman build -f Dockerfile-go -t <registry_endpoint>/tulospalvelu:latest .
# podman login <registry_endpoint> -u nologin --password-stdin <<< "$SCW_SECRET_KEY"
# podman push <registry_endpoint>/tulospalvelu:latest
```

### 3. Deploy the container

```bash
cd terraform
terraform apply -var deploy_container=true
```

### 4. Open the app

```bash
terraform output container_url   # add -json if terraform complains about sensitivity
```

Wait for the container status to be `ready`:

```bash
scw container container list name=tulospalvelu -o human=Name,Status
```

Then open the URL. Note the first request may be slow (cold start of
container and database — by design, both scale to zero).

## Operations

- **New version**: rebuild + push the image (same tag `latest`, or a
  new tag via `container_image_tag`), then `terraform apply` again.
- **DB connection** (psql): `terraform output db_connection_string`.
- **Teardown**: `terraform destroy` (removes everything, including all
  data in the Serverless SQL Database).

## Notes

- The app migrates the schema automatically on startup
  (`CREATE TABLE IF NOT EXISTS`). No initdb scripts are needed.
- Write operations (POST) are restricted to Finnish IP addresses
  (`GEO_RESTRICT=true` in the container env). The client IP is taken
  from the rightmost `X-Forwarded-For` entry (`TRUST_XFF=true`), which
  Scaleway's edge sets. GET endpoints and the UI are public.
- Database auth uses IAM: the DB username is the IAM application ID,
  the password is its API key secret — both managed by Terraform and
  injected as the container's secret `DATABASE_URL` env var
  (`sslmode=require`).
