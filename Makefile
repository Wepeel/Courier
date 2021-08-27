all: build-casey docker-casey
build-casey:
	go build -ldflags "-s -w" -o bin/casey cmd/Casey/main.go
docker-casey:
	docker build --pull --rm -f ./build/docker/Casey/Dockerfile -t casey:latest .
build-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/app/protos/*.proto