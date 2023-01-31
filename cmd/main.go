package main

import (
	"fmt"
	"time"

	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap"

	"github.com/saltpay/transaction-card-validator/internal/adapters/http"
	"github.com/saltpay/transaction-card-validator/internal/adapters/prometheus"
	"github.com/saltpay/transaction-card-validator/internal/application/ports"
)

func main() {
	ctx, done := listenForCancellationAndAddToContext()
	defer done()

	service, err := newService(ctx)
	if err != nil {
		zapctx.Fatal(ctx, "failed to create service", zap.Error(err))
	}

	go service.Listener.Listen(ctx)

	zapctx.Info(ctx, "Bootstrapped all the components. Starting the various services now...")
	serverConfig := http.ServerConfig{
		Port:             "8080",
		HTTPReadTimeout:  2 * time.Second,
		HTTPWriteTimeout: 2 * time.Second,
	}

	server := http.NewWebServer(serverConfig)

	zapctx.Info(ctx, fmt.Sprintf("Started. Listening on port: %s", serverConfig.Port))
	if err = server.ListenAndServe(); err != nil {
		zapctx.Fatal(ctx, "http server listen failed", zap.Error(err))
	}
}

func newMetricsClient() ports.MetricsClient {
	prometheusClient := prometheus.New()

	return &prometheusClient
}
