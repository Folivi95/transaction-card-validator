package main

import (
	"strings"

	zapctx "github.com/saltpay/go-zap-ctx"

	logger "github.com/saltpay/transaction-card-validator/internal/adapters/prometheus/zap"

	"go.uber.org/zap/zapcore"
)

func InitLog(appConfig AppConfig) {
	switch strings.ToLower(appConfig.LogConfig.LogLevel) {
	case "debug":
		zapctx.SetLogLevel(zapcore.DebugLevel)
	case "info":
		zapctx.SetLogLevel(zapcore.InfoLevel)
	case "error":
		zapctx.SetLogLevel(zapcore.ErrorLevel)
	case "warn":
		zapctx.SetLogLevel(zapcore.WarnLevel)
	default:
		zapctx.SetLogLevel(zapcore.DebugLevel)
	}

	// Add error log hook (for alerting)
	hook := logger.NewPrometheusZapHook()
	zapctx.AddHooks(hook)
}
