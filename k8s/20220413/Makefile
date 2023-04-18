include .env.local
export

APP1_DOCKER_REPO := $(DOCKERHUB_USERNAME)/app1
APP2_DOCKER_REPO := $(DOCKERHUB_USERNAME)/app2

app1:
	@go run cmd/app1/main.go

app2:
	@go run cmd/app2/main.go

image_app1:
	@KO_DOCKER_REPO=$(APP1_DOCKER_REPO) ko build --sbom=none --bare ./cmd/app1 --platform=linux/amd64

image_app2:
	@KO_DOCKER_REPO=$(APP2_DOCKER_REPO) ko build --sbom=none --bare ./cmd/app2 --platform=linux/amd64

update:
	@kubectl delete -f k8s.yaml
	@kubectl apply -f k8s.yaml 

logs_app1:
	@kubectl logs -l app=app1 -c app1

logs_app2:
	@kubectl logs -l app=app2 -c app2

port_forward_app2:
	@kubectl port-forward svc/app2 8081:8080
          
.PHONY: gen
gen:
	@go generate ./...
	@(cd proto && buf generate --template buf.gen.yaml)
	@go mod tidy