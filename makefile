apigen:
	@echo "Running oapi-codegen..."
	@oapi-codegen --config=api/config.yaml -o internal/api/generated/server.gen.go -package api api/openapi.yaml
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run
