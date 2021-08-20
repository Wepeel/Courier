FROM golang:1.16-alpine AS build-env
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/main ./cmd/Casey/main.go

FROM scratch
WORKDIR /app
COPY --from=build-env /app/bin/main /app/
CMD ["./main"]