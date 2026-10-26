.PHONY: fmt vet test run image compose-up compose-down tf-init tf-fmt tf-validate tf-plan tf-apply-bootstrap tf-apply tf-destroy

fmt:
	gofmt -l .

vet:
	go vet ./...

test:
	go test ./...

run:
	DATABASE_URL=postgres://tulospalvelu:tulospalvelu@localhost:5432/tulospalvelu?sslmode=disable \
	SEED_TEST_DATA=true go run ./cmd/server

compose-up:
	podman compose -f compose.go.yaml up --build

compose-down:
	podman compose -f compose.go.yaml down -v

# Build the production image. Usage:
#   make image REGISTRY=$(cd terraform && terraform output -raw registry_image)
image:
	podman build -f Dockerfile-go -t $(REGISTRY) .

tf-init:
	cd terraform && terraform init

tf-fmt:
	cd terraform && terraform fmt

tf-validate:
	cd terraform && terraform validate

tf-plan:
	cd terraform && terraform plan

tf-apply-bootstrap:
	cd terraform && terraform apply

tf-apply:
	cd terraform && terraform apply -var deploy_container=true

tf-destroy:
	cd terraform && terraform destroy
