build-casey:
	go build -ldflags "-s -w" -o bin/casey cmd/Casey/main.go
docker-casey:
	docker build --pull --rm -f ./build/docker/Casey/Dockerfile -t casey:latest .
all: build-casey docker-casey