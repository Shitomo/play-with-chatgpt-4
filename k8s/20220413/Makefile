app1:
	@go run cmd/app1/main.go

app2:
	@go run cmd/app2/main.go

image_app1:
	@ko build --sbom=none --bare  --local ./cmd/app1 --platform=linux/amd64