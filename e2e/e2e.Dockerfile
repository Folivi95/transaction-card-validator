FROM golang:1.18 as deps

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY vendor ./vendor
COPY go.mod go.sum ./

COPY e2e/ ./e2e/
COPY scripts/ ./scripts/

CMD [ "go", "test", "-count=1", "--tags=e2e", "./..." ]
