# Generate Go gRPC code from proto (run from backend/)
# Requires: protoc, protoc-gen-go, protoc-gen-go-grpc on PATH
.PHONY: proto-go
proto-go:
	protoc --go_out=go-backend --go_opt=paths=source_relative \
		--go-grpc_out=go-backend --go-grpc_opt=paths=source_relative \
		--proto_path=. \
		proto/ats/screening.proto
	@echo "Go code in go-backend/proto/ats/"

# Generate Python gRPC code (run from backend/ai-service with venv active)
.PHONY: proto-python
proto-python:
	cd ai-service && python -m grpc_tools.protoc -I.. \
		--python_out=app/grpc_gen --grpc_python_out=app/grpc_gen \
		../proto/ats/screening.proto
	@echo "Python code in ai-service/app/grpc_gen/proto/ats/"
