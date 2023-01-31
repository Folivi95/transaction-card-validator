FROM public.ecr.aws/docker/library/golang:1.17.2 as deps

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPRIVATE=github.com/saltpay/*

RUN git config --global url."git@github.com".insteadOf "https://github.com"

WORKDIR /app
COPY vendor ./vendor
COPY go.mod go.sum ./

COPY cmd/ ./cmd/
COPY internal/ ./internal/

COPY certs/salt.pem /usr/share/ca-certificates/
RUN sed -i -e '$asalt.pem' /etc/ca-certificates.conf
RUN update-ca-certificates

FROM deps as build
RUN go build -o transaction-card-validator github.com/saltpay/transaction-card-validator/cmd

FROM public.ecr.aws/docker/library/alpine:3.10
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/ ./

CMD ["/app/transaction-card-validator"]
